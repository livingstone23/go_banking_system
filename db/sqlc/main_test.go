package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

const (
	dbDriver = "pgx"
	dbSource = "postgres://alumno:123456@localhost:5432/simple_bank?sslmode=disable"
)

// this var is used to store the queries
var testQueries *Queries

// This function is the main function of the test
// It is used to connect to the database
func TestMain(m *testing.M) {
	conn, err := pgx.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	defer conn.Close(context.Background())

	testQueries = New(conn)

	os.Exit(m.Run())
}
