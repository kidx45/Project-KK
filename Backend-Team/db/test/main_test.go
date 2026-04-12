package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	sqlc "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"

)

const driverName = "postgres"
const dataSourceName = "postgresql://root:secret@localhost:5433/Project_KK?sslmode=disable"

var testQueries *sqlc.Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(driverName, dataSourceName)

	if err != nil {
		log.Fatal("cannot connect to the database:", err)
	}

	testQueries = sqlc.New(conn)
	exitCode := m.Run()
	conn.Close()
	os.Exit(exitCode)
}
