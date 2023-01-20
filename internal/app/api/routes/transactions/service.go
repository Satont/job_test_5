package transactions

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/satont/test/internal/amqp"
	"github.com/satont/test/internal/app/api"
	"github.com/satont/test/internal/constants"
	"github.com/satont/test/internal/db/models"
	"github.com/satont/test/internal/dto"
	"github.com/shopspring/decimal"
)

type ReplenishValidationDto struct {
	Amount decimal.Decimal `validate:"required" json:"amount"`
}

type WithDrawValidationDto struct {
	Amount decimal.Decimal `validate:"required" json:"amount"`
}

var NotEnoughBalance = errors.New("not enough balance")

func GetService(app *api.App, consumerId, transactionId string) (*models.Transaction, error) {
	transaction, err := app.TransactionsService.GetById(transactionId, consumerId)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func GetManyService(app *api.App, consumerId string, limit, offset uint64, status models.TransactionStatus) ([]models.Transaction, error) {
	transactions, err := app.TransactionsService.GetManyByConsumerId(consumerId, limit, offset, status)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func ReplenishService(app *api.App, consumerId string, validationDto *ReplenishValidationDto) error {
	transactionDto := dto.Transaction{
		Amount:     validationDto.Amount,
		ConsumerID: consumerId,
		Type:       models.TypeReplenish,
	}

	bytes, err := json.Marshal(&transactionDto)
	if err != nil {
		return err
	}

	app.TransactionsAmqpChannel.PublishWithContext(context.Background(),
		"",
		constants.WITHDRAW_CHANNEL_QUEUE,
		false,
		false,
		amqp.CreateMessage(bytes),
	)

	return nil
}

func WithDrawService(app *api.App, consumerId string, validationDto *WithDrawValidationDto) error {
	consumerBalance, err := app.ConsumersService.GetBalanceById(consumerId)
	if err != nil {
		app.Logger.Error(err)
		return err
	}

	transactionDto := dto.Transaction{
		Amount:     validationDto.Amount,
		ConsumerID: consumerId,
		Type:       models.TypeWithDraw,
	}

	if consumerBalance.LessThan(transactionDto.Amount) {
		return NotEnoughBalance
	}

	bytes, err := json.Marshal(&transactionDto)
	if err != nil {
		return err
	}

	app.TransactionsAmqpChannel.PublishWithContext(context.Background(),
		"",
		constants.WITHDRAW_CHANNEL_QUEUE,
		false,
		false,
		amqp.CreateMessage(bytes),
	)

	return nil
}
