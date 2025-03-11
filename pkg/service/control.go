package service

import (
	"github.com/Diyarjan/BankSystem"
	"github.com/Diyarjan/BankSystem/pkg/repository"
)

type ControlService struct {
	repo repository.Control
}

func NewControlService(repo repository.Control) *ControlService {
	return &ControlService{repo: repo}
}

func (s *ControlService) CreateAccount(account BankSystem.ToMakeAccount) (int, error) {
	return s.repo.CreateAccount(account)
}

func (s *ControlService) DeleteAccount(accountID int) error {
	return s.repo.DeleteAccount(accountID)
}

func (s *ControlService) GetAccountByID(accountID int) (BankSystem.Account, error) {
	return s.repo.GetAccountByID(accountID)
}

func (s *ControlService) GetAllAccounts() ([]BankSystem.Account, error) {
	return s.repo.GetAllAccounts()
}
