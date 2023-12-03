package model

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	TransactionDebit  string = "debit"
	TransactionCredit string = "credit"
)

type AccountTransaction struct {
	ID              uint64          `json:"id" db:"id"`
	AccountID       uint64          `json:"account_id" db:"account_id"`
	TransactionType string          `json:"transaction_type" db:"transaction_type"`
	Remark          string          `json:"remark" db:"remark"`
	Amount          decimal.Decimal `json:"balance" db:"balance"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
}

type NewTransaction struct {
	UserID          uint64          `json:"user_id"`
	Amount          decimal.Decimal `json:"amount"`
	Remark          string          `json:"remark"`
	TransactionType string          `json:"transaction_type"`
	ReferenceKey    string          `json:"reference_key"`
}
