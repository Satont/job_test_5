package consumers_impl

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/satont/test/internal/db/models"
	"github.com/shopspring/decimal"
)

type ConsumersService struct {
	pgConn *sqlx.DB
}

func NewSqlxConsumers(pgConn *sqlx.DB) *ConsumersService {
	return &ConsumersService{pgConn: pgConn}
}

func (c *ConsumersService) GetById(id string) (*models.Consumer, error) {
	query, args, err := sq.
		Select("*").
		From("consumers").
		Where(sq.Eq{"id": id}).
		ToSql()
	query = c.pgConn.Rebind(query)

	if err != nil {
		return nil, err
	}

	consumer := &models.Consumer{}
	err = c.pgConn.Get(consumer, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return consumer, nil
}

func (c *ConsumersService) GetByApiKey(apiKey string) (*models.Consumer, error) {
	query, args, err := sq.
		Select("*").
		From("consumers").
		Where(sq.Eq{"api_key": apiKey}).
		ToSql()
	query = c.pgConn.Rebind(query)

	if err != nil {
		return nil, err
	}

	consumer := &models.Consumer{}
	err = c.pgConn.Get(consumer, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return consumer, nil
}

func (c *ConsumersService) GetBalanceById(consumerId string) (*decimal.Decimal, error) {
	query, args, err := sq.
		Select().
		Column(`SUM(
			CASE WHEN type = ? AND status != ? THEN -amount
			WHEN type = ? AND status = ? THEN amount
			ELSE 0 END
		) AS balance`, models.TypeWithDraw, models.StatusCanceled, models.TypeReplenish, models.StatusProcessed).
		From("transactions").
		Where(sq.Eq{"consumer_id": consumerId}).
		GroupBy("consumer_id").
		ToSql()

	if err != nil {
		return nil, err
	}
	query = c.pgConn.Rebind(query)

	rows, err := c.pgConn.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	balance := decimal.Decimal{}

	for rows.Next() {
		var b decimal.Decimal
		err = rows.Scan(&b)
		if err != nil {
			return nil, err
		}
		balance = b
	}

	return &balance, nil
}
