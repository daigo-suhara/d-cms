package repository

import (
	"context"

	"github.com/daigo-suhara/d-cms/internal/domain"
)

type ContentModelRepository interface {
	FindAll(ctx context.Context) ([]domain.ContentModel, error)
	FindByID(ctx context.Context, id uint) (*domain.ContentModel, error)
	FindBySlug(ctx context.Context, slug string) (*domain.ContentModel, error)
	Create(ctx context.Context, m *domain.ContentModel) error
	Update(ctx context.Context, m *domain.ContentModel) error
	Delete(ctx context.Context, id uint) error
}

type EntryRepository interface {
	FindAllByModel(ctx context.Context, modelID uint) ([]domain.Entry, error)
	FindByID(ctx context.Context, id uint) (*domain.Entry, error)
	Create(ctx context.Context, e *domain.Entry) error
	Update(ctx context.Context, e *domain.Entry) error
	Delete(ctx context.Context, id uint) error
	CountByModel(ctx context.Context, modelID uint) (int64, error)
}

type MediaRepository interface {
	FindAll(ctx context.Context) ([]domain.Media, error)
	FindByID(ctx context.Context, id uint) (*domain.Media, error)
	Create(ctx context.Context, m *domain.Media) error
	Delete(ctx context.Context, id uint) error
}
