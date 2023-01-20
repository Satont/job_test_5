package transactions_impl

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/satont/test/internal/db/models"
	"github.com/satont/test/internal/services/transactions"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"time"
)

type TransactionsService struct {
	pgConn *sqlx.DB
}

func NewSqlxTransactions(pgConn *sqlx.DB) transactions.TransactionsService {
	return &TransactionsService{
		pgConn: pgConn,
	}
}

func (c *TransactionsService) GetById(transactionId, consumerId string) (*models.Transaction, error) {
	query, args, err := sq.
		Select("*").
		From("transactions").
		Where(sq.Eq{"id": transactionId, "consumer_id": consumerId}).
		ToSql()
	query = c.pgConn.Rebind(query)

	if err != nil {
		return nil, err
	}

	transaction := &models.Transaction{}
	err = c.pgConn.Get(transaction, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return transaction, nil
}

func (c *TransactionsService) GetManyByConsumerId(
	consumerId string,
	limit,
	offset uint64,
	status models.TransactionStatus,
) ([]models.Transaction, error) {
	queryBuilder := sq.
		Select("*").
		From("transactions").
		Where(sq.Eq{"consumer_id": consumerId}).
		OrderBy("created_at DESC").
		Limit(limit).
		Offset(offset)

	if status != "" {
		queryBuilder = queryBuilder.Where(sq.Eq{"status": status})
	}

	query, args, err := queryBuilder.ToSql()
	query = c.pgConn.Rebind(query)

	if err != nil {
		return nil, err
	}

	transactions := []models.Transaction{}
	err = c.pgConn.Select(&transactions, query, args...)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (c *TransactionsService) Create(
	consumerId string,
	amount decimal.Decimal,
	t models.TransactionType,
) (*models.Transaction, error) {
	amountValue, err := amount.Value()
	if err != nil {
		return nil, err
	}

	transactionId := uuid.NewV4().String()
	query, args, err := sq.
		Insert("transactions").
		Columns("id", "amount", "status", "type", "consumer_id", "created_at").
		Values(
			transactionId,
			amountValue,
			models.StatusCreated,
			t,
			consumerId,
			time.Now().UTC(),
		).
		ToSql()
	query = c.pgConn.Rebind(query)

	if err != nil {
		return nil, err
	}

	_, err = c.pgConn.Queryx(query, args...)

	if err != nil {
		return nil, err
	}

	transaction, err := c.GetById(transactionId, consumerId)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
