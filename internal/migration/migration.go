package migration

import (
	"fmt"

	"github.com/daigo-suhara/d-cms/internal/domain"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&domain.ContentModel{},
		&domain.Entry{},
		&domain.Media{},
		&domain.APIKey{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}
