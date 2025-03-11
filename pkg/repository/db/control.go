package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Diyarjan/BankSystem"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"time"
)

type ControlPostgres struct {
	db    *sqlx.DB
	cache *redis.Client
}

var ctx = context.Background()

func NewControlPostgres(db *sqlx.DB, cache *redis.Client) *ControlPostgres {
	return &ControlPostgres{db: db, cache: cache}
}

func (r *ControlPostgres) CreateAccount(account BankSystem.ToMakeAccount) (int, error) {
	var newAccount BankSystem.Account

	query := "INSERT INTO accounts(balance, currency) VALUES ($1, $2) RETURNING id, balance, currency, is_locked, created_at;"
	if err := r.db.Get(&newAccount, query, account.Balance, account.Currency); err != nil {
		fmt.Println("Error inserting new account")
		return 0, err
	}

	// Insert Account data into Redis with expire time of 24 hours
	accountKey := fmt.Sprintf("account-%d", newAccount.Id)
	accountJson, _ := json.Marshal(newAccount)
	if err := r.cache.SetEx(ctx, accountKey, accountJson, 24*time.Hour).Err(); err != nil {
		fmt.Println("CreateAccount redis set ex error:", err)
	}

	return newAccount.Id, nil
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

	// Delete Account data from Redis
	accountKey := fmt.Sprintf("account-%d", accountID)
	err = r.cache.Del(context.Background(), accountKey).Err()
	if err != nil {
		return fmt.Errorf("failed to cache account data in Redis: %s", err)
	}

	return nil
}

func (r *ControlPostgres) GetAccountByID(accountID int) (BankSystem.Account, error) {
	var account BankSystem.Account

	// Get from Redis
	accountKey := fmt.Sprintf("account-%d", accountID)
	resp, _ := r.cache.Get(ctx, accountKey).Bytes()
	if resp != nil {
		_ = json.Unmarshal(resp, &account)
	}

	if account.Id > 0 {
		return account, nil
	}

	// Get from DB
	query := "SELECT id, balance, currency, is_locked, created_at FROM accounts WHERE id = $1;"
	if err := r.db.Get(&account, query, accountID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return BankSystem.Account{}, errors.New("account id = %d does not exist")
		}
		return BankSystem.Account{}, err
	}

	// If exist in Postgres then Insert to Redis
	accountJson, _ := json.Marshal(account)
	if err := r.cache.SetEx(context.Background(), accountKey, accountJson, 24*time.Hour).Err(); err != nil {
		fmt.Println("GetByID redis set ex error:", err)
	}
	fmt.Println("exist in postgres and insert into Redis")

	return account, nil
}

func (r *ControlPostgres) GetAllAccounts() ([]BankSystem.Account, error) {
	var accounts []BankSystem.Account

	query := "SELECT id, balance, currency, is_locked, created_at FROM accounts;"
	if err := r.db.Select(&accounts, query); err != nil {
		return []BankSystem.Account{}, err
	}

	return accounts, nil
}
