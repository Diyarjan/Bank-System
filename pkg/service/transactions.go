package service

import (
	"encoding/json"
	"fmt"
	"github.com/Diyarjan/BankSystem"
	"github.com/Diyarjan/BankSystem/constants"
	"github.com/Diyarjan/BankSystem/pkg/repository"
	"github.com/Diyarjan/BankSystem/third_party/kafkaPart"
	"log"
)

type TransactionService struct {
	repo     repository.Transaction
	producer *kafkaPart.Producer
}

func NewTransactionService(repo repository.Transaction, producer *kafkaPart.Producer) *TransactionService {
	return &TransactionService{repo: repo, producer: producer}
}

func (s *TransactionService) DepositToAccount(depositStruct BankSystem.DebitCreditStruct) error {
	_, err := s.repo.CheckValidateAccount(depositStruct.AccountID)
	if err != nil {
		return err
	}
	if depositStruct.Amount <= 0 {
		log.Println("Deposit amount is less than zero")
		return fmt.Errorf("invalid amount: %f ", depositStruct.Amount)
	}

	message, err := json.Marshal(depositStruct)
	if err != nil {
		log.Printf("Field marshaled from service %s", err)
		return err
	}

	if err := s.producer.SendMessage([]byte(message), constants.Deposit); err != nil {
		log.Printf("Field send message %s", string(message))
		return err
	}
	return nil
	//return s.repo.DepositToAccount(accountID, amount)
}

func (s *TransactionService) WithdrawFromAccount(withdrawStruct BankSystem.DebitCreditStruct) error {
	account, err := s.repo.CheckValidateAccount(withdrawStruct.AccountID)
	if err != nil {
		return err
	}

	if withdrawStruct.Amount <= 0 {
		log.Println("Withdraw amount is less than zero")
		return fmt.Errorf("invalid amount: %f ", withdrawStruct.Amount)
	}

	if account.Balance < withdrawStruct.Amount {
		log.Println("Withdraw amount is big from your balance", withdrawStruct.Amount)
		return fmt.Errorf("Withdraw amount is big from your balance: %f ", account.Balance)
	}

	message, err := json.Marshal(withdrawStruct)
	if err != nil {
		log.Printf("Field marshaled from service %s", err)
		return err
	}
	if err := s.producer.SendMessage([]byte(message), constants.Withdraw); err != nil {
		log.Printf("Field send message %s", string(message))
		return err
	}
	return nil
	//return s.repo.WithdrawFromAccount(accountID, amount)
}

func (s *TransactionService) TransferFunds(transferStruct BankSystem.Transfer) (float32, error) {
	// Check validation from account
	account, err := s.repo.CheckValidateAccount(transferStruct.FromAccountId)
	if err != nil {
		return 0, err
	}
	if transferStruct.Amount <= 0 {
		log.Println("Transfer amount is less than zero")
		return 0, fmt.Errorf("Transfer amount is less than zero: %f ", transferStruct.Amount)
	}
	if account.Balance < transferStruct.Amount {
		log.Println("Transfer amount is big from your balance", account.Balance)
		return 0, fmt.Errorf("Transfer amount is big from your balance: %f ", account.Balance)
	}
	oldBalance := account.Balance
	// Check validation to account
	account, err = s.repo.CheckValidateAccount(transferStruct.ToAccountID)
	if err != nil {
		return 0, err
	}

	message, err := json.Marshal(transferStruct)
	if err != nil {
		log.Printf("Field marshaled from service %s", err)
		return 0, err
	}

	if err := s.producer.SendMessage([]byte(message), constants.Transfer); err != nil {
		log.Printf("Field send message %s", string(message))
		return 0, err
	}
	return oldBalance - transferStruct.Amount, nil
	//return s.repo.TransferFunds(accountID, transferStruct)
}

func (s *TransactionService) GetTransactionHistory(accountID int) ([]BankSystem.Transaction, error) {
	return s.repo.GetTransactionHistory(accountID)
}
