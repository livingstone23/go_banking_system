package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver = "pgx"
	dbSource = "postgres://alumno:123456@localhost:5432/simple_bank?sslmode=disable"
)

// this var is used to store the queries
var testQueries *Queries

// this var is used to store the connection to the database
//var testDB *pgx.Conn
var testDB *pgxpool.Pool




func TestMain(m *testing.M) {
    var err error
    testDB, err = pgxpool.New(context.Background(), dbSource)
    if err != nil {
        log.Fatal("cannot connect to db:", err)
    }
    defer testDB.Close()

    testQueries = New(testDB)

    // Run the tests
    code := m.Run()

    // Close the pool after tests are done
    testDB.Close()

    // Exit with the test code
    os.Exit(code)
}

/*

func TestMain(m *testing.M) {
    var err error
    testDB, err = pgxpool.Connect(context.Background(), dbSource)
    if err != nil {
        log.Fatal("cannot connect to db:", err)
    }
    defer testDB.Close()

    testQueries = New(testDB)

    os.Exit(m.Run())
}
// This function is the main function of the test
// It is used to connect to the database
func TestMain(m *testing.M) {

	var err error
	testDB, err = pgx.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// this is used to close the connection to the database
	//defer testDB.Close(context.Background())

	testQueries = New(testDB)

	os.Exit(m.Run())
}
	*/
