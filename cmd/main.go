package main

import (
	"github.com/delapaska/EgorExplain/pkg/postgres"
	"github.com/delapaska/EgorExplain/pkg/router"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := postgres.NewPostgresSQLStorage()
	if err != nil {
		return
	}
	if err = postgres.RunMigrations(db, "migrations"); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	server := router.NewApiServer(db)

	server.Run()
}
