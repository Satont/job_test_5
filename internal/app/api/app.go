package api

import (
	"github.com/rabbitmq/amqp091-go"
	"github.com/satont/test/internal/services/consumers"
	"github.com/satont/test/internal/services/logger"
	"github.com/satont/test/internal/services/transactions"
)

type App struct {
	ConsumersService    consumers.ConsumersService
	TransactionsService transactions.TransactionsService
	Logger              logger.Logger

	TransactionsAmqpChannel *amqp091.Channel
}
