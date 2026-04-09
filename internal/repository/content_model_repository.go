package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/daigo-suhara/d-cms/internal/domain"
	"gorm.io/gorm"
)

type contentModelRepository struct {
	db *gorm.DB
}

func NewContentModelRepository(db *gorm.DB) ContentModelRepository {
	return &contentModelRepository{db: db}
}

func (r *contentModelRepository) FindAll(ctx context.Context) ([]domain.ContentModel, error) {
	var models []domain.ContentModel
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("FindAll content models: %w", err)
	}
	return models, nil
}

func (r *contentModelRepository) FindByID(ctx context.Context, id uint) (*domain.ContentModel, error) {
	var m domain.ContentModel
	err := r.db.WithContext(ctx).First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID content model %d: %w", id, err)
	}
	return &m, nil
}

func (r *contentModelRepository) FindBySlug(ctx context.Context, slug string) (*domain.ContentModel, error) {
	var m domain.ContentModel
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("FindBySlug content model %q: %w", slug, err)
	}
	return &m, nil
}

func (r *contentModelRepository) Create(ctx context.Context, m *domain.ContentModel) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("Create content model: %w", err)
	}
	return nil
}

func (r *contentModelRepository) Update(ctx context.Context, m *domain.ContentModel) error {
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("Update content model %d: %w", m.ID, err)
	}
	return nil
}

func (r *contentModelRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&domain.ContentModel{}, id).Error; err != nil {
		return fmt.Errorf("Delete content model %d: %w", id, err)
	}
	return nil
}
