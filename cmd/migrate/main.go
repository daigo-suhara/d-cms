package main

import (
	"log"

	"github.com/daigo-suhara/d-cms/config"
	"github.com/daigo-suhara/d-cms/internal/infrastructure/database"
	"github.com/daigo-suhara/d-cms/internal/migration"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Running database migrations...")
	if err := migration.Run(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations completed successfully.")
}
