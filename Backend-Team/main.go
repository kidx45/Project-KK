package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kidx45/Project-KK/Backend-Team/api"
	db "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	"github.com/kidx45/Project-KK/Backend-Team/utils"
	_ "github.com/lib/pq"
)

func main() {
	AppConfig, err := utils.LoadEnv("app.env")
	if err != nil {
		log.Fatal("Can't load data because: ", err)
	}

	conn, err := sql.Open(AppConfig.DB_DRIVER_NAME, AppConfig.DB_URL)
	if err != nil {
		log.Fatal("Can't start server due to: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(AppConfig, store)
	if err != nil {
		log.Fatal("Can't create server due to: ", err)
	}
	address := fmt.Sprintf("localhost:%s", AppConfig.PORT)
	err = server.Start(address)
	if err != nil {
		log.Fatal("Can't start server due to: ", err)
	}
}
