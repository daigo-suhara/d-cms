package repository

import (
	"context"
	"errors"

	"github.com/daigo-suhara/d-cms/internal/domain"
	"gorm.io/gorm"
)

type apiKeyRepository struct{ db *gorm.DB }

func NewAPIKeyRepository(db *gorm.DB) APIKeyRepository {
	return &apiKeyRepository{db: db}
}

func (r *apiKeyRepository) FindAll(ctx context.Context) ([]domain.APIKey, error) {
	var keys []domain.APIKey
	if err := r.db.WithContext(ctx).Order("created_at desc").Find(&keys).Error; err != nil {
		return nil, err
	}
	return keys, nil
}

func (r *apiKeyRepository) FindByKey(ctx context.Context, key string) (*domain.APIKey, error) {
	var k domain.APIKey
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&k).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	return &k, err
}

func (r *apiKeyRepository) Create(ctx context.Context, k *domain.APIKey) error {
	return r.db.WithContext(ctx).Create(k).Error
}

func (r *apiKeyRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.APIKey{}, id).Error
}
