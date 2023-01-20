package consumers

import (
	"github.com/satont/test/internal/db/models"
	"github.com/shopspring/decimal"
)

type ConsumersService interface {
	GetById(consumerId string) (*models.Consumer, error)
	GetByApiKey(apiKey string) (*models.Consumer, error)
	GetBalanceById(consumerId string) (*decimal.Decimal, error)
}
