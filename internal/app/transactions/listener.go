package transactions

import (
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/satont/test/internal/dto"
	"github.com/satont/test/internal/services/logger"
	"github.com/satont/test/internal/services/transactions"
)

func NewEventsListener(ch <-chan amqp091.Delivery, service transactions.TransactionsService, logger logger.Logger) {
	go func() {
		for d := range ch {
			logger.Info(fmt.Sprintf("Received a message: %s", d.Body))

			transactionDto := dto.Transaction{}
			err := json.Unmarshal(d.Body, &transactionDto)
			if err != nil {
				logger.Error(err)
				return
			}
			_, err = service.Create(transactionDto.ConsumerID, transactionDto.Amount, transactionDto.Type)
			if err != nil {
				d.Ack(false)
			}
		}
	}()
}
