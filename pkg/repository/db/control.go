package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Diyarjan/BankSystem"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type ControlPostgres struct {
	db    *sqlx.DB
	cache *redis.Client
}

var ctx = context.Background()

func NewControlPostgres(db *sqlx.DB, cache *redis.Client) *ControlPostgres {
	return &ControlPostgres{db: db, cache: cache}
}

func (r *ControlPostgres) CreateAccount(account BankSystem.ToMakeAccount) (BankSystem.Account, error) {
	var newAccount BankSystem.Account
	query := "INSERT INTO accounts(balance, currency) VALUES ($1, $2) RETURNING id, balance, currency, is_locked, created_at;"
	if err := r.db.Get(&newAccount, query, account.Balance, account.Currency); err != nil {
		fmt.Println("Error inserting new account")
		return BankSystem.Account{}, err
	}
	return newAccount, nil
}

func (r *ControlPostgres) DeleteAccount(accountID int) error {
	result, err := r.db.Exec("UPDATE accounts SET is_locked = true, deleted_at = now() WHERE id = $1;", accountID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("account %d does not exist", accountID)
	}

	fmt.Println("Account deleted from Postgres", accountID)
	return nil
}

func (r *ControlPostgres) GetAccountByID(accountID int) (BankSystem.Account, error) {
	var account BankSystem.Account
	query := "SELECT id, balance, currency, is_locked, created_at FROM accounts WHERE id = $1;"
	if err := r.db.Get(&account, query, accountID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return BankSystem.Account{}, errors.New("account not exist")
		}
		return BankSystem.Account{}, err
	}
	return account, nil
}

func (r *ControlPostgres) GetAllAccounts() ([]BankSystem.Account, error) {
	var accounts []BankSystem.Account

	query := "SELECT id, balance, currency, is_locked, created_at FROM accounts WHERE is_locked = false;"
	if err := r.db.Select(&accounts, query); err != nil {
		return []BankSystem.Account{}, err
	}

	return accounts, nil
}
