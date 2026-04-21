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
	config, err := utils.LoadEnv()
	if err != nil {
		log.Fatal("cannot load env:", err)
	}

	testDB, err = sql.Open(config.DB_DRIVER_NAME, config.DB_URL)

	if err != nil {
		log.Fatal("cannot connect to the database:", err)
	}

	testQueries = sqlc.New(testDB)
	exitCode := m.Run()
	testDB.Close()
	os.Exit(exitCode)
}
