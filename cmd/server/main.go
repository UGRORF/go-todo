package main

import (
	"fmt"
	"os"

	"github.com/UGRORF/go-todo/internal/api"
	"github.com/UGRORF/go-todo/pkg/db"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	dbFile := os.Getenv("TODO_DBFILE")

	err := db.Init(dbFile)
	if err != nil {
		fmt.Println(err)
	}

	port := os.Getenv("TODO_PORT")
	srv := api.NewServer(port)

	if err := srv.Start(); err != nil {
		fmt.Printf("Server failed: %v", err)
	}

}
