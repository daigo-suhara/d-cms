package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/daigo-suhara/d-cms/internal/domain"
	"github.com/daigo-suhara/d-cms/internal/repository"
)

var slugRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

type ContentModelService struct {
	modelRepo repository.ContentModelRepository
	entryRepo repository.EntryRepository
}

func NewContentModelService(modelRepo repository.ContentModelRepository, entryRepo repository.EntryRepository) *ContentModelService {
	return &ContentModelService{modelRepo: modelRepo, entryRepo: entryRepo}
}

func (s *ContentModelService) List(ctx context.Context) ([]domain.ContentModel, error) {
	return s.modelRepo.FindAll(ctx)
}

func (s *ContentModelService) GetByID(ctx context.Context, id uint) (*domain.ContentModel, error) {
	return s.modelRepo.FindByID(ctx, id)
}

func (s *ContentModelService) GetBySlug(ctx context.Context, slug string) (*domain.ContentModel, error) {
	return s.modelRepo.FindBySlug(ctx, slug)
}

func (s *ContentModelService) Create(ctx context.Context, m *domain.ContentModel) error {
	if err := s.validateModel(m); err != nil {
		return err
	}
	existing, err := s.modelRepo.FindBySlug(ctx, m.Slug)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("check slug uniqueness: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("slug %q: %w", m.Slug, domain.ErrAlreadyExists)
	}
	if m.Fields == nil {
		m.Fields = []domain.FieldDefinition{}
	}
	return s.modelRepo.Create(ctx, m)
}

func (s *ContentModelService) Update(ctx context.Context, m *domain.ContentModel) error {
	if err := s.validateModel(m); err != nil {
		return err
	}
	existing, err := s.modelRepo.FindBySlug(ctx, m.Slug)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("check slug uniqueness: %w", err)
	}
	if existing != nil && existing.ID != m.ID {
		return fmt.Errorf("slug %q: %w", m.Slug, domain.ErrAlreadyExists)
	}
	if m.Fields == nil {
		m.Fields = []domain.FieldDefinition{}
	}
	return s.modelRepo.Update(ctx, m)
}

func (s *ContentModelService) Delete(ctx context.Context, id uint) error {
	count, err := s.entryRepo.CountByModel(ctx, id)
	if err != nil {
		return fmt.Errorf("check entries before delete: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("cannot delete model with %d existing entries: %w", count, domain.ErrHasEntries)
	}
	return s.modelRepo.Delete(ctx, id)
}

func (s *ContentModelService) validateModel(m *domain.ContentModel) error {
	if m.Name == "" {
		return fmt.Errorf("name is required: %w", domain.ErrValidation)
	}
	if m.Slug == "" {
		return fmt.Errorf("slug is required: %w", domain.ErrValidation)
	}
	if !slugRegex.MatchString(m.Slug) {
		return fmt.Errorf("slug must match ^[a-z0-9-]+$: %w", domain.ErrValidation)
	}
	seen := map[string]bool{}
	for _, f := range m.Fields {
		if f.Name == "" {
			return fmt.Errorf("field name is required: %w", domain.ErrValidation)
		}
		if seen[f.Name] {
			return fmt.Errorf("duplicate field name %q: %w", f.Name, domain.ErrValidation)
		}
		seen[f.Name] = true
		switch f.Type {
		case domain.FieldTypeText, domain.FieldTypeNumber, domain.FieldTypeDate, domain.FieldTypeMarkdown:
		default:
			return fmt.Errorf("unknown field type %q: %w", f.Type, domain.ErrValidation)
		}
	}
	return nil
}
