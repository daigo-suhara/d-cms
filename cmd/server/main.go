package main

import (
	"log"

	"github.com/daigo-suhara/d-cms/config"
	"github.com/daigo-suhara/d-cms/internal/infrastructure/database"
	"github.com/daigo-suhara/d-cms/internal/infrastructure/storage"
	"github.com/daigo-suhara/d-cms/internal/migration"
	"github.com/daigo-suhara/d-cms/router"
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

	r2, err := storage.NewR2Client(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage client: %v", err)
	}

	r := router.Setup(db, r2, cfg)
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
