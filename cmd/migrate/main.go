package main

import (
	"log"

	"github.com/daigo-suhara/d-cms/config"
	"github.com/daigo-suhara/d-cms/internal/domain"
	"github.com/daigo-suhara/d-cms/internal/infrastructure/database"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Running database migrations...")
	if err := db.AutoMigrate(
		&domain.ContentModel{},
		&domain.Entry{},
		&domain.Media{},
		&domain.APIKey{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations completed successfully.")
}
