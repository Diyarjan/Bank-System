package service

import (
	"github.com/Diyarjan/BankSystem"
	"github.com/Diyarjan/BankSystem/pkg/repository"
	"github.com/Diyarjan/BankSystem/third_party/kafkaPart"
)

type Control interface {
	CreateAccount(account BankSystem.ToMakeAccount) (int, error)
	DeleteAccount(accountID int) error
	GetAccountByID(accountID int) (BankSystem.Account, error)
	GetAllAccounts() ([]BankSystem.Account, error)
}

type Transaction interface {
	DepositToAccount(depositStruct BankSystem.DebitCreditStruct) error
	WithdrawFromAccount(withdrawStruct BankSystem.DebitCreditStruct) error
	TransferFunds(transfer BankSystem.Transfer) (float32, error)
	GetTransactionHistory(accountID int) ([]BankSystem.Transaction, error)
}

type Service struct {
	Control
	Transaction
}

func NewService(repos *repository.Repository, producer *kafkaPart.Producer) *Service {
	return &Service{
		Control:     NewControlService(repos.Control),
		Transaction: NewTransactionService(repos.Transaction, producer),
	}
}
