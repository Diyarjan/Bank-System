package BankSystem

import (
	"time"
)

type ToMakeAccount struct {
	Balance  int    `json:"balance"`
	Currency string `json:"currency" binding:"required"`
}

type Account struct {
	Id        int       `json:"id" db:"id"`
	Balance   float32   `json:"balance" db:"balance"`
	Currency  string    `json:"currency" db:"currency"`
	IsLocked  bool      `json:"is_locked" db:"is_locked"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	//DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"`
}

type Transaction struct {
	AccountId       int       `json:"account_id" db:"account_id"`
	Amount          float32   `json:"amount" db:"amount"`
	TransactionType string    `json:"transaction_type" db:"transaction_type"`
	CreatedAt       time.Time `json:"transaction_time" db:"created_at"`
}

type Transfer struct {
	FromAccountId int     `json:"from_account_id"`
	ToAccountID   int     `json:"to_account_id" binding:"required"`
	Amount        float32 `json:"amount" binding:"required"`
}

type DebitCreditStruct struct {
	AccountID int     `json:"account_id" binding:"required"`
	Amount    float32 `json:"amount" binding:"required"`
}

type CheckValidationStruct struct {
	ID       int     `json:"id" db:"id"`
	Balance  float32 `json:"balance" db:"balance"`
	IsLocked bool    `json:"is_locked" db:"is_locked"`
}
