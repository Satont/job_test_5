package dto

import (
	"github.com/satont/test/internal/db/models"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	Amount     decimal.Decimal        `json:"amount"`
	ConsumerID string                 `json:"consumer_id"`
	Type       models.TransactionType `json:"type"`
}
