package repository

import (
	"github.com/Diyarjan/BankSystem"
	"github.com/Diyarjan/BankSystem/pkg/repository/db"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Control interface {
	CreateAccount(account BankSystem.ToMakeAccount) (BankSystem.Account, error)
	DeleteAccount(accountID int) error
	GetAccountByID(accountID int) (BankSystem.Account, error)
	GetAllAccounts() ([]BankSystem.Account, error)
}

type Transaction interface {
	DepositToAccount(depositStruct BankSystem.DebitCreditStruct) error
	WithdrawFromAccount(withdrawStruct BankSystem.DebitCreditStruct) error
	TransferFunds(transfer BankSystem.Transfer) (float32, error)
	GetTransactionHistory(accountID int) ([]BankSystem.Transaction, error)

	CheckValidateAccount(id int) (BankSystem.CheckValidationStruct, error)
}
type Repository struct {
	Control
	Transaction
}

func NewRepository(dbConn *sqlx.DB, cache *redis.Client) *Repository {
	return &Repository{
		Control:     db.NewControlPostgres(dbConn, cache),
		Transaction: db.NewTransactionPostgres(dbConn, cache),
	}
}
