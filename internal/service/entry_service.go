package service

import (
	"context"
	"fmt"
	"time"

	"github.com/daigo-suhara/d-cms/internal/domain"
	"github.com/daigo-suhara/d-cms/internal/repository"
)

type EntryService struct {
	modelRepo repository.ContentModelRepository
	entryRepo repository.EntryRepository
}

func NewEntryService(modelRepo repository.ContentModelRepository, entryRepo repository.EntryRepository) *EntryService {
	return &EntryService{modelRepo: modelRepo, entryRepo: entryRepo}
}

func (s *EntryService) ListBySlug(ctx context.Context, modelSlug string) ([]domain.Entry, *domain.ContentModel, error) {
	model, err := s.modelRepo.FindBySlug(ctx, modelSlug)
	if err != nil {
		return nil, nil, err
	}
	entries, err := s.entryRepo.FindAllByModel(ctx, model.ID)
	if err != nil {
		return nil, nil, err
	}
	return entries, model, nil
}

func (s *EntryService) GetByID(ctx context.Context, id uint) (*domain.Entry, error) {
	return s.entryRepo.FindByID(ctx, id)
}

func (s *EntryService) Create(ctx context.Context, modelSlug string, content domain.ContentData) (*domain.Entry, error) {
	model, err := s.modelRepo.FindBySlug(ctx, modelSlug)
	if err != nil {
		return nil, err
	}
	if err := s.validate(model, content); err != nil {
		return nil, err
	}
	e := &domain.Entry{
		ContentModelID: model.ID,
		Content:        content,
	}
	if err := s.entryRepo.Create(ctx, e); err != nil {
		return nil, err
	}
	e.ContentModel = *model
	return e, nil
}

func (s *EntryService) Update(ctx context.Context, id uint, content domain.ContentData) (*domain.Entry, error) {
	e, err := s.entryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.validate(&e.ContentModel, content); err != nil {
		return nil, err
	}
	e.Content = content
	if err := s.entryRepo.Update(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *EntryService) Delete(ctx context.Context, id uint) error {
	if _, err := s.entryRepo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.entryRepo.Delete(ctx, id)
}

func (s *EntryService) validate(model *domain.ContentModel, content domain.ContentData) error {
	for _, field := range model.Fields {
		val, exists := content[field.Name]
		isEmpty := !exists || val == nil || val == ""

		if field.Required && isEmpty {
			return fmt.Errorf("field %q is required: %w", field.Name, domain.ErrValidation)
		}
		if isEmpty {
			continue
		}

		switch field.Type {
		case domain.FieldTypeNumber:
			if _, ok := val.(float64); !ok {
				return fmt.Errorf("field %q must be a number: %w", field.Name, domain.ErrInvalidField)
			}
		case domain.FieldTypeDate:
			str, ok := val.(string)
			if !ok {
				return fmt.Errorf("field %q must be an ISO 8601 date string: %w", field.Name, domain.ErrInvalidField)
			}
			if _, err := time.Parse(time.RFC3339, str); err != nil {
				if _, err2 := time.Parse("2006-01-02T15:04", str); err2 != nil {
					return fmt.Errorf("field %q must be ISO 8601 (e.g. 2006-01-02T15:04:05Z): %w", field.Name, domain.ErrInvalidField)
				}
			}
		case domain.FieldTypeText, domain.FieldTypeMarkdown:
			if _, ok := val.(string); !ok {
				return fmt.Errorf("field %q must be a string: %w", field.Name, domain.ErrInvalidField)
			}
		case domain.FieldTypeTags:
			items, ok := val.([]any)
			if !ok {
				return fmt.Errorf("field %q must be an array of strings: %w", field.Name, domain.ErrInvalidField)
			}
			for _, item := range items {
				if _, ok := item.(string); !ok {
					return fmt.Errorf("field %q must be an array of strings: %w", field.Name, domain.ErrInvalidField)
				}
			}
		}
	}
	return nil
}
