package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/paulohhs/ms-balance/internal/database"
	"github.com/paulohhs/ms-balance/internal/usecase/get_balance"
	"github.com/paulohhs/ms-balance/internal/usecase/update_balance"
	"github.com/paulohhs/ms-balance/internal/web"
	"github.com/paulohhs/ms-balance/internal/web/webserver"
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
		return database.NewAccountDB(tx) // tx, não db: as queries rodam dentro da transação
	})

	updateBalanceUseCase := update_balance.NewUpdateBalanceUseCase(unitOfWork)
	_ = updateBalanceUseCase

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
