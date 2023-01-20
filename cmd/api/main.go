package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"github.com/satont/test/internal/amqp"
	"github.com/satont/test/internal/app/api"
	"github.com/satont/test/internal/app/api/router"
	consumers_impl "github.com/satont/test/internal/services/consumers/impl"
	logger_impl "github.com/satont/test/internal/services/logger/impl"
	transactions_impl "github.com/satont/test/internal/services/transactions/impl"
	"log"
	"net/http"
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

	app := &api.App{
		ConsumersService:    consumers_impl.NewSqlxConsumers(db),
		TransactionsService: transactions_impl.NewSqlxTransactions(db),
		Logger:              logger_impl.NewConsoleLogger(),

		TransactionsAmqpChannel: ch,
	}

	router := router.InitRoutes(app)
	go http.ListenAndServe("0.0.0.0:8080", router)

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	fmt.Println("Closing...")
}
