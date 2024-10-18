package main

import (
	"context"
	"fmt"
	"go_banking_system/api"
	db "go_banking_system/db/sqlc"
	"log"
	"go_banking_system/util"
	"github.com/jackc/pgx/v5/pgxpool"
)



func main() {
	fmt.Println("Hello, from Banking System!!!!")


	// Load parameters from the environment
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	defer conn.Close()

	fmt.Println("Connected to the database")

	// Create a new server
	store := db.NewStore(conn)
	server := api.NewServer(store)

	// Start the server
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
