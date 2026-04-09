package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/daigo-suhara/d-cms/internal/domain"
	"gorm.io/gorm"
)

type mediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) MediaRepository {
	return &mediaRepository{db: db}
}

func (r *mediaRepository) FindAll(ctx context.Context) ([]domain.Media, error) {
	var media []domain.Media
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&media).Error; err != nil {
		return nil, fmt.Errorf("FindAll media: %w", err)
	}
	return media, nil
}

func (r *mediaRepository) FindByID(ctx context.Context, id uint) (*domain.Media, error) {
	var m domain.Media
	err := r.db.WithContext(ctx).First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID media %d: %w", id, err)
	}
	return &m, nil
}

func (r *mediaRepository) Create(ctx context.Context, m *domain.Media) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("Create media: %w", err)
	}
	return nil
}

func (r *mediaRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Media{}, id).Error; err != nil {
		return fmt.Errorf("Delete media %d: %w", id, err)
	}
	return nil
}
