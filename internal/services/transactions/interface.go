package transactions

import (
	"github.com/satont/test/internal/db/models"
	"github.com/shopspring/decimal"
)

type TransactionStatus string

type TransactionsService interface {
	GetById(id, consumerId string) (*models.Transaction, error)
	GetManyByConsumerId(consumerId string, limit, offset uint64, status models.TransactionStatus) ([]models.Transaction, error)
	Create(consumerId string, amount decimal.Decimal, t models.TransactionType) (*models.Transaction, error)
}
