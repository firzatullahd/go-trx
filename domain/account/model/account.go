package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID        uint64          `json:"id" db:"id"`
	UserID    uint64          `json:"user_id" db:"user_id"`
	Balance   decimal.Decimal `json:"balance" db:"balance"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time      `json:"deleted_at" db:"deleted_at"`
}
