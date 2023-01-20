package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"github.com/satont/test/internal/amqp"
	logger_impl "github.com/satont/test/internal/services/logger/impl"
	transactions_impl "github.com/satont/test/internal/services/transactions/impl"

	"github.com/satont/test/internal/app/transactions"
	"github.com/satont/test/internal/constants"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	godotenv.Load()
}

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

	rabbitUrl := os.Getenv("RABBIT_URL")
	conn, err := amqp.CreateAmqpConn(rabbitUrl)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	ch, err := amqp.CreateChannel(conn)
	if err != nil {
		log.Panic(err)
	}
	defer ch.Close()

	q, err := amqp.CreateQueue(ch, constants.WITHDRAW_CHANNEL_QUEUE)
	if err != nil {
		log.Panic(err)
	}

	msgs, err := amqp.CreateConsume(ch, q.Name)
	if err != nil {
		log.Panic(err)
	}

	logger := logger_impl.NewConsoleLogger()
	transactionsService := transactions_impl.NewSqlxTransactions(db)

	transactions.NewEventsListener(msgs, transactionsService, logger)

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	fmt.Println("Closing...")
}
