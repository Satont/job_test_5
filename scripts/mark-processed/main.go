package main

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"os"
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

	query, args, err := sq.
		Update("transactions").
		Set("status", "processed").
		ToSql()
	query = db.Rebind(query)

	if err != nil {
		panic(err)
	}

	_, err = db.Queryx(query, args...)

	if err != nil {
		panic(err)
	}

	if _, err = db.Query(query, args...); err != nil {
		panic(err)
	}
}
