package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	sqlc "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	utils "github.com/kidx45/Project-KK/Backend-Team/utils"
)

var testQueries *sqlc.Queries
var testDB *sql.DB
func TestMain(m *testing.M) {
	config := utils.Config{
		DB_URL:         os.Getenv("DB_URL"),
		PORT:           os.Getenv("PORT"),
		DB_DRIVER_NAME: os.Getenv("DB_DRIVER_NAME"),
	}

	if config.DB_URL == "" || config.DB_DRIVER_NAME == "" {
		log.Fatal("DB_URL or DB_DRIVER_NAME is not set in environment variables")
	}

	var err error
	testDB, err = sql.Open(config.DB_DRIVER_NAME, config.DB_URL)

	if err != nil {
		log.Fatal("cannot connect to the database:", err)
	}

	testQueries = sqlc.New(testDB)
	exitCode := m.Run()
	testDB.Close()
	os.Exit(exitCode)
}
