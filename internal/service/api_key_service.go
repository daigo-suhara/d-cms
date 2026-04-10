package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/daigo-suhara/d-cms/internal/domain"
	"github.com/daigo-suhara/d-cms/internal/repository"
)

type APIKeyService struct {
	repo repository.APIKeyRepository
}

func NewAPIKeyService(repo repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{repo: repo}
}

func (s *APIKeyService) List(ctx context.Context) ([]domain.APIKey, error) {
	return s.repo.FindAll(ctx)
}

func (s *APIKeyService) Create(ctx context.Context, name string) (*domain.APIKey, error) {
	if name == "" {
		return nil, fmt.Errorf("名前は必須です: %w", domain.ErrValidation)
	}
	key, err := generateKey()
	if err != nil {
		return nil, fmt.Errorf("キー生成に失敗しました: %w", err)
	}
	k := &domain.APIKey{Name: name, Key: key}
	if err := s.repo.Create(ctx, k); err != nil {
		return nil, err
	}
	return k, nil
}

func (s *APIKeyService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *APIKeyService) Validate(ctx context.Context, key string) bool {
	_, err := s.repo.FindByKey(ctx, key)
	return err == nil
}

func generateKey() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "dcms_" + hex.EncodeToString(b), nil
}
