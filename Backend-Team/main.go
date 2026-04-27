package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kidx45/Project-KK/Backend-Team/api"
	db "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load(".env.local.development")

	DB_URL := os.Getenv("DB_URL")
	DB_DRIVER_NAME := os.Getenv("DB_DRIVER_NAME")
	PORT := os.Getenv("PORT")

	conn, err := sql.Open(DB_DRIVER_NAME, DB_URL)
	if err != nil {
		log.Fatal("Can't start seerver due to: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	address := fmt.Sprintf("localhost:%s", PORT)
	err = server.Start(address)
	if err != nil {
		log.Fatal("Can't start server due to: ", err)
	}
}
