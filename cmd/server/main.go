package main

import (
	"fmt"
	"os"

	"github.com/UGRORF/go-todo/internal/api"
	"github.com/UGRORF/go-todo/pkg/db"
	"github.com/joho/godotenv"
)

const (
	defaultPort   = "7540"
	defaultDBFile = "./scheduler.db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}

	err := db.Init(dbFile)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}
	srv := api.NewServer(port)

	if err := srv.Start(); err != nil {
		fmt.Printf("Server failed: %v", err)
	}

}
