package main

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/satont/test/internal/db/models"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
	"time"
)

func main() {
	dbUrl := os.Getenv("POSTGRES_URL")
	dbConnOpts, err := pq.ParseURL(dbUrl)
	if err != nil {
		panic(fmt.Errorf("cannot parse postgres url connection: %w", err))
	}
	db, err := sqlx.Connect("postgres", dbConnOpts)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	apiKey := "5e3e90b4-ec51-4f94-8550-902e3cfa5d6a"

	query, args, err := sq.
		Select("*").
		From("consumers").
		Where(sq.Eq{"api_key": apiKey}).
		ToSql()
	query = db.Rebind(query)

	if err != nil {
		panic(err)
	}

	consumer := &models.Consumer{}
	err = db.Get(consumer, query, args...)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	if consumer.ID != "" {
		return
	}

	query, args, err = sq.
		Insert("consumers").
		Columns("id", "name", "created_at", "api_key").
		Values(
			uuid.NewV4().String(),
			"test-consumer",
			time.Now().UTC(),
			apiKey,
		).
		ToSql()
	query = db.Rebind(query)

	if err != nil {
		panic(err)
	}

	if _, err = db.Query(query, args...); err != nil {
		panic(err)
	}
}
