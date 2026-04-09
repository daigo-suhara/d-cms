package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/daigo-suhara/d-cms/internal/domain"
	"gorm.io/gorm"
)

type entryRepository struct {
	db *gorm.DB
}

func NewEntryRepository(db *gorm.DB) EntryRepository {
	return &entryRepository{db: db}
}

func (r *entryRepository) FindAllByModel(ctx context.Context, modelID uint) ([]domain.Entry, error) {
	var entries []domain.Entry
	err := r.db.WithContext(ctx).
		Where("content_model_id = ?", modelID).
		Order("created_at DESC").
		Find(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("FindAllByModel entries for model %d: %w", modelID, err)
	}
	return entries, nil
}

func (r *entryRepository) FindByID(ctx context.Context, id uint) (*domain.Entry, error) {
	var e domain.Entry
	err := r.db.WithContext(ctx).Preload("ContentModel").First(&e, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID entry %d: %w", id, err)
	}
	return &e, nil
}

func (r *entryRepository) Create(ctx context.Context, e *domain.Entry) error {
	if err := r.db.WithContext(ctx).Create(e).Error; err != nil {
		return fmt.Errorf("Create entry: %w", err)
	}
	return nil
}

func (r *entryRepository) Update(ctx context.Context, e *domain.Entry) error {
	if err := r.db.WithContext(ctx).Save(e).Error; err != nil {
		return fmt.Errorf("Update entry %d: %w", e.ID, err)
	}
	return nil
}

func (r *entryRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Entry{}, id).Error; err != nil {
		return fmt.Errorf("Delete entry %d: %w", id, err)
	}
	return nil
}

func (r *entryRepository) CountByModel(ctx context.Context, modelID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Entry{}).
		Where("content_model_id = ?", modelID).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("CountByModel entries for model %d: %w", modelID, err)
	}
	return count, nil
}
