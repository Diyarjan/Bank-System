package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Diyarjan/BankSystem"
	"github.com/Diyarjan/BankSystem/constants"
	"github.com/redis/go-redis/v9"
	"time"
)

var ctx = context.Background()

type RedisService interface {
	CreateAccount(account BankSystem.Account) error
	DeleteAccount(accountID int) error
	GetAccountByID(accountID int) (BankSystem.Account, error)
	GetAllAccounts() ([]BankSystem.Account, error)
	SetToRedis(accountKey string, account []byte) error
}

type redisService struct {
	cache *redis.Client
}

func NewRedisService(cacheClient *redis.Client) RedisService {
	return &redisService{cache: cacheClient}
}

func (s *redisService) CreateAccount(account BankSystem.Account) error {
	accountKey := fmt.Sprintf("account-%d", account.Id)
	accountJson, _ := json.Marshal(account)
	if err := s.cache.SetEx(ctx, accountKey, accountJson, 24*time.Hour).Err(); err != nil {
		fmt.Println("CreateAccount redis set ex error:", err)
		return err
	}

	// update list id account
	var accountIDs []int
	cachedIDs, err := s.cache.Get(ctx, constants.IdListKey).Bytes()
	if err == nil && len(cachedIDs) > 0 {
		_ = json.Unmarshal(cachedIDs, &accountIDs)
	}

	accountIDs = append(accountIDs, account.Id)
	idListJSON, _ := json.Marshal(accountIDs)
	_ = s.cache.Set(ctx, constants.IdListKey, idListJSON, 10*time.Minute).Err()

	return nil
}

func (s *redisService) DeleteAccount(accountID int) error {
	accountKey := fmt.Sprintf("account-%d", accountID)
	if err := s.cache.Del(ctx, accountKey).Err(); err != nil {
		//return fmt.Errorf("%s Error delete account from redis", err)
		return err
	}
	fmt.Println("Account not exist in redis for delete", accountID)
	return nil
}

func (s *redisService) GetAccountByID(accountID int) (BankSystem.Account, error) {
	var account BankSystem.Account
	accountKey := fmt.Sprintf("account-%d", accountID)
	resp, err := s.cache.Get(ctx, accountKey).Bytes()
	if err != nil {
		return BankSystem.Account{}, err
	}

	if err := json.Unmarshal(resp, &account); err != nil {
		return BankSystem.Account{}, err
	}

	if account.Id > 0 {
		fmt.Println("From cachePart:", account)
		return account, nil
	}

	return BankSystem.Account{}, fmt.Errorf(fmt.Sprintf("account %d not found", accountID))
}

func (s *redisService) GetAllAccounts() ([]BankSystem.Account, error) {

	// 1. Получаем список всех ID аккаунтов
	cachedIDs, err := s.cache.Get(ctx, constants.IdListKey).Bytes()
	var accountIDs []int
	if err == nil && len(cachedIDs) > 0 {
		_ = json.Unmarshal(cachedIDs, &accountIDs)
	}
	// 2. Загружаем аккаунты из кеша
	var accounts []BankSystem.Account
	var account BankSystem.Account
	for _, accountID := range accountIDs {
		accountKey := fmt.Sprintf("account-%d", accountID)
		cachedAccount, err := s.cache.Get(ctx, accountKey).Bytes()
		if err == nil {
			account = BankSystem.Account{}
			_ = json.Unmarshal(cachedAccount, &account)
			accounts = append(accounts, account)
		}
	}

	// 3. Если в кеше все данные есть – возвращаем их
	if len(accounts) == len(accountIDs) && len(accounts) > 0 {
		return accounts, nil
	}

	return accounts, fmt.Errorf("some accounts not found")
}

func (s *redisService) SetToRedis(accountKey string, account []byte) error {
	return s.cache.Set(ctx, accountKey, account, 10*time.Minute).Err()
}
