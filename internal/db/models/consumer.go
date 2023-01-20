package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type Consumer struct {
	ID        string    `db:"id" json:"id"`
	ApiKey    string    `db:"api_key" json:"-"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	Balance *decimal.Decimal `json:"-"`
}
