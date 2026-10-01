package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/paulohhs/ms-balance/internal/database"
	"github.com/paulohhs/ms-balance/internal/event"
	"github.com/paulohhs/ms-balance/internal/event/handler"
	"github.com/paulohhs/ms-balance/internal/usecase/get_balance"
	"github.com/paulohhs/ms-balance/internal/usecase/update_balance"
	"github.com/paulohhs/ms-balance/internal/web"
	"github.com/paulohhs/ms-balance/internal/web/webserver"
	"github.com/paulohhs/ms-balance/pkg/events"
	"github.com/paulohhs/ms-balance/pkg/kafka"
	"github.com/paulohhs/ms-balance/pkg/uow"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		getEnv("DB_USER", "root"),
		getEnv("DB_PASSWORD", "root"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "3307"),
		getEnv("DB_NAME", "balance"),
	))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	accountDb := database.NewAccountDB(db)
	getBalanceUseCase := get_balance.NewGetBalanceUseCase(accountDb)

	ctx := context.Background()
	unitOfWork := uow.NewUow(ctx, db)
	unitOfWork.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(tx) // tx, as queries rodam dentro da transação
	})
	updateBalanceUseCase := update_balance.NewUpdateBalanceUseCase(unitOfWork)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("BalanceUpdated", handler.NewBalanceUpdatedKafkaHandler(updateBalanceUseCase))

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
		"group.id":          "balances",
		"auto.offset.reset": "earliest",
	}
	consumer := kafka.NewConsumer(&configMap, []string{"balances"})
	msgChan := make(chan *ckafka.Message)
	go consumer.Consume(msgChan)
	go func() {
		for msg := range msgChan {
			balanceUpdated := event.NewBalanceUpdated()
			if err := json.Unmarshal(msg.Value, balanceUpdated); err != nil {
				fmt.Println("error decoding kafka message:", err)
				continue
			}
			eventDispatcher.Dispatch(balanceUpdated)
		}
	}()

	server := webserver.NewWebServer(":3003")
	balanceHandler := web.NewWebBalanceHandler(*getBalanceUseCase)
	server.AddHandler(http.MethodGet, "/balances/{account_id}", balanceHandler.GetBalance)

	fmt.Println("Server is running on port 3003")
	if err := server.Start(); err != nil {
		panic(err)
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
