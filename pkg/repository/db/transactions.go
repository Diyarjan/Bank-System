package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Diyarjan/BankSystem"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"log"
)

type TransactionPostgres struct {
	db    *sqlx.DB
	cache *redis.Client
}

func NewTransactionPostgres(db *sqlx.DB, cache *redis.Client) *TransactionPostgres {
	return &TransactionPostgres{db: db, cache: cache}
}

func (r *TransactionPostgres) DepositToAccount(depositStruct BankSystem.DebitCreditStruct) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	var balance float32
	query := "UPDATE accounts SET balance = accounts.balance + $1 WHERE id = $2 AND is_locked = false RETURNING balance"
	row := tx.QueryRow(query, depositStruct.Amount, depositStruct.AccountID)
	if err := row.Scan(&balance); err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("not found")
		}
		return err
	}

	//	Insert to transaction table
	query = `INSERT INTO transactions (account_id, amount, transaction_type) VALUES ($1, $2, $3);`
	_, err = tx.Exec(query, depositStruct.AccountID, depositStruct.Amount, "deposit")
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *TransactionPostgres) WithdrawFromAccount(withdrawStruct BankSystem.DebitCreditStruct) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var balance float32
	query := "SELECT balance FROM accounts WHERE id = $1 AND is_locked = false"
	row := tx.QueryRow(query, withdrawStruct.AccountID)
	if err := row.Scan(&balance); err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("not found")
		}
		return err
	}
	if balance < withdrawStruct.Amount {
		tx.Rollback()
		return errors.New("not enough funds")
	}

	query = "UPDATE accounts SET balance = accounts.balance - $1 WHERE id = $2 AND is_locked = false RETURNING balance"
	row = tx.QueryRow(query, withdrawStruct.Amount, withdrawStruct.AccountID)
	if err := row.Scan(&balance); err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("not found")
		}
		return err
	}

	//	Insert to transaction table
	query = `INSERT INTO transactions (account_id, amount, transaction_type) VALUES ($1, $2, $3);`
	_, err = tx.Exec(query, withdrawStruct.AccountID, withdrawStruct.Amount, "credit")
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *TransactionPostgres) TransferFunds(transfer BankSystem.Transfer) (float32, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var balance float32
	var id int
	//	Find balance
	query := "Select balance FROM accounts WHERE id = $1 AND is_locked = false"
	row := r.db.QueryRow(query, transfer.FromAccountId)
	if err := row.Scan(&balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("not found")
		}
		return 0, err
	}
	if balance < transfer.Amount {
		return balance, errors.New("not enough funds")
	}

	//	Withdraw
	query = "UPDATE accounts SET balance = accounts.balance - $1 WHERE id = $2 AND is_locked = false RETURNING balance"
	row = tx.QueryRow(query, transfer.Amount, transfer.FromAccountId)
	if err = row.Scan(&balance); err != nil {
		tx.Rollback()
		return 0, err
	}

	//	Deposit
	query = "UPDATE accounts SET balance = accounts.balance + $1 WHERE id = $2 AND is_locked = false RETURNING id"
	row = tx.QueryRow(query, transfer.Amount, transfer.ToAccountID)
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			return 0, errors.New("not found")
		}
		_ = tx.Rollback()
		return 0, err
	}

	//	Insert to transactions table
	query = `INSERT INTO transactions (account_id, amount, transaction_type) VALUES ($1, $2, $3);`
	_, err = tx.Exec(query, transfer.FromAccountId, transfer.Amount, "credit")
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	query = `INSERT INTO transactions (account_id, amount, transaction_type) VALUES ($1, $2, $3);`
	_, err = tx.Exec(query, transfer.ToAccountID, transfer.Amount, "deposit")
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return balance, tx.Commit()
}

func (r *TransactionPostgres) GetTransactionHistory(accountID int) ([]BankSystem.Transaction, error) {
	var transactions []BankSystem.Transaction
	query := "SELECT account_id, amount, transaction_type, created_at FROM transactions WHERE account_id = $1 AND deleted_at is NULL ORDER BY created_at DESC"
	err := r.db.Select(&transactions, query, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []BankSystem.Transaction{}, errors.New("not found")
		}
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionPostgres) CheckValidateAccount(id int) (BankSystem.CheckValidationStruct, error) {
	var account BankSystem.CheckValidationStruct
	query := "SELECT id, balance, is_locked FROM accounts WHERE id = $1"
	if err := r.db.Get(&account, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("account %d does not exist", id)
			return BankSystem.CheckValidationStruct{}, fmt.Errorf("account %d does not exist", id)
		}
	}
	if account.IsLocked {
		log.Printf("account %d is locked", id)
		return BankSystem.CheckValidationStruct{}, fmt.Errorf("account %d is locked", id)
	}
	return account, nil
}

//IsLockedAccount(id int) error
//CheckBalance(id int) (float64, error)
