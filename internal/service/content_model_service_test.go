package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/daigo-suhara/d-cms/internal/domain"
	"github.com/daigo-suhara/d-cms/internal/service"
)

func newCMSvc(models map[string]*domain.ContentModel, entries []*domain.Entry) *service.ContentModelService {
	cmRepo := &mockContentModelRepo{models: models}
	entryRepo := &mockEntryRepo{entries: entries}
	return service.NewContentModelService(cmRepo, entryRepo)
}

func TestContentModelService_Create_Valid(t *testing.T) {
	svc := newCMSvc(map[string]*domain.ContentModel{}, nil)
	m := &domain.ContentModel{
		Name:   "Article",
		Slug:   "article",
		Fields: []domain.FieldDefinition{{Name: "title", Type: domain.FieldTypeText, Required: true}},
	}
	if err := svc.Create(context.Background(), m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestContentModelService_Create_InvalidSlug(t *testing.T) {
	svc := newCMSvc(map[string]*domain.ContentModel{}, nil)
	m := &domain.ContentModel{Name: "Test", Slug: "Invalid Slug!"}
	err := svc.Create(context.Background(), m)
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestContentModelService_Create_DuplicateSlug(t *testing.T) {
	existing := &domain.ContentModel{ID: 1, Name: "Existing", Slug: "blog"}
	svc := newCMSvc(map[string]*domain.ContentModel{"blog": existing}, nil)
	m := &domain.ContentModel{Name: "Blog2", Slug: "blog"}
	err := svc.Create(context.Background(), m)
	if !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestContentModelService_Delete_WithEntries(t *testing.T) {
	existing := &domain.ContentModel{ID: 1, Name: "Blog", Slug: "blog"}
	entries := []*domain.Entry{{ID: 1, ContentModelID: 1}}
	svc := newCMSvc(map[string]*domain.ContentModel{"blog": existing}, entries)
	err := svc.Delete(context.Background(), 1)
	if !errors.Is(err, domain.ErrHasEntries) {
		t.Fatalf("expected ErrHasEntries, got %v", err)
	}
}

func TestContentModelService_Delete_Empty(t *testing.T) {
	existing := &domain.ContentModel{ID: 1, Name: "Blog", Slug: "blog"}
	svc := newCMSvc(map[string]*domain.ContentModel{"blog": existing}, nil)
	if err := svc.Delete(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestContentModelService_Create_DuplicateFieldName(t *testing.T) {
	svc := newCMSvc(map[string]*domain.ContentModel{}, nil)
	m := &domain.ContentModel{
		Name: "Test",
		Slug: "test",
		Fields: []domain.FieldDefinition{
			{Name: "title", Type: domain.FieldTypeText},
			{Name: "title", Type: domain.FieldTypeMarkdown},
		},
	}
	err := svc.Create(context.Background(), m)
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected ErrValidation for duplicate field, got %v", err)
	}
}
