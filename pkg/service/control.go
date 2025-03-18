package service

import (
	"encoding/json"
	"fmt"
	"github.com/Diyarjan/BankSystem"
	"github.com/Diyarjan/BankSystem/constants"
	"github.com/Diyarjan/BankSystem/pkg/repository"
	"github.com/Diyarjan/BankSystem/pkg/service/cache"
	"github.com/redis/go-redis/v9"
)

type ControlService struct {
	repo  repository.Control
	cache cache.RedisService
}

func NewControlService(repo repository.Control, rdsClient *redis.Client) *ControlService {
	return &ControlService{
		repo:  repo,
		cache: cache.NewRedisService(rdsClient),
	}
}

func (s *ControlService) CreateAccount(account BankSystem.ToMakeAccount) (int, error) {

	//Create account in Postgres
	newAccount, err := s.repo.CreateAccount(account)
	if err != nil {
		fmt.Println("Postgres createAccount err:", err)
		return 0, err
	}

	//Insert new account into Redis
	if err = s.cache.CreateAccount(newAccount); err != nil {
		fmt.Println("cache create account error", err)
		return newAccount.Id, err
	}

	fmt.Println("Postgres & Redis createAccount success")
	return newAccount.Id, nil
}

func (s *ControlService) DeleteAccount(accountID int) error {
	if err := s.cache.DeleteAccount(accountID); err != nil {
		return err
	}
	return s.repo.DeleteAccount(accountID)
}

func (s *ControlService) GetAccountByID(accountID int) (BankSystem.Account, error) {
	//Get by cache if exist
	account, err := s.cache.GetAccountByID(accountID)
	if err == nil && account.Id > 0 {
		fmt.Println("Return from cache")
		return account, err
	}

	//Get by postgres and insert into cache if not exist cache,
	account, err = s.repo.GetAccountByID(accountID)
	if err != nil {
		return BankSystem.Account{}, err
	}
	if err = s.cache.CreateAccount(account); err != nil {
		fmt.Printf("cache create account err:%v", err)
	}

	return account, nil
}

func (s *ControlService) GetAllAccounts() ([]BankSystem.Account, error) {

	accounts, err := s.cache.GetAllAccounts()
	if err == nil {
		fmt.Println("Accounts return from cache")
		return accounts, nil
	}

	accounts, err = s.repo.GetAllAccounts()
	if err != nil {
		return accounts, err
	}

	// Обновляем кеш: сохраняем ID + аккаунты по отдельности
	var accountIDs []int
	for _, account := range accounts {
		accountIDs = append(accountIDs, account.Id)
		accountKey := fmt.Sprintf("account-%d", account.Id)
		jsonData, _ := json.Marshal(account)
		_ = s.cache.SetToRedis(accountKey, jsonData)
	}

	// Записываем список всех ID в Redis
	idListJSON, _ := json.Marshal(accountIDs)
	_ = s.cache.SetToRedis(constants.IdListKey, idListJSON)

	fmt.Println("Accounts returning after update from postgres")
	return accounts, nil
}
