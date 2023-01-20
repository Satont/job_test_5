package models

import (
	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v4"
	"time"
)

type TransactionStatus string

const (
	StatusCreated    TransactionStatus = "created"
	StatusProcessing TransactionStatus = "processing"
	StatusProcessed  TransactionStatus = "processed"
	StatusCanceled   TransactionStatus = "canceled"
)

type TransactionType string

const (
	TypeWithDraw  TransactionType = "withdraw"
	TypeReplenish TransactionType = "replenish"
)

type Transaction struct {
	ID         string            `db:"id" json:"id"`
	Amount     decimal.Decimal   `db:"amount" json:"amount"`
	Status     TransactionStatus `db:"status" json:"status"`
	Type       TransactionType   `db:"type" json:"type"`
	ConsumerID string            `db:"consumer_id" json:"consumer_id"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt null.Time `db:"updated_at" json:"updated_at"`
}
