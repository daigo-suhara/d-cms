package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/daigo-suhara/d-cms/internal/domain"
	"github.com/daigo-suhara/d-cms/internal/service"
)

// --- Minimal in-memory mocks ---

type mockContentModelRepo struct {
	models map[string]*domain.ContentModel
}

func (m *mockContentModelRepo) FindAll(_ context.Context) ([]domain.ContentModel, error) {
	var out []domain.ContentModel
	for _, v := range m.models {
		out = append(out, *v)
	}
	return out, nil
}
func (m *mockContentModelRepo) FindByID(_ context.Context, id uint) (*domain.ContentModel, error) {
	for _, v := range m.models {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, domain.ErrNotFound
}
func (m *mockContentModelRepo) FindBySlug(_ context.Context, slug string) (*domain.ContentModel, error) {
	v, ok := m.models[slug]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return v, nil
}
func (m *mockContentModelRepo) Create(_ context.Context, cm *domain.ContentModel) error {
	m.models[cm.Slug] = cm
	return nil
}
func (m *mockContentModelRepo) Update(_ context.Context, cm *domain.ContentModel) error {
	m.models[cm.Slug] = cm
	return nil
}
func (m *mockContentModelRepo) Delete(_ context.Context, id uint) error {
	for k, v := range m.models {
		if v.ID == id {
			delete(m.models, k)
			return nil
		}
	}
	return domain.ErrNotFound
}

type mockEntryRepo struct {
	entries []*domain.Entry
	nextID  uint
}

func (m *mockEntryRepo) FindAllByModel(_ context.Context, modelID uint) ([]domain.Entry, error) {
	var out []domain.Entry
	for _, e := range m.entries {
		if e.ContentModelID == modelID {
			out = append(out, *e)
		}
	}
	return out, nil
}
func (m *mockEntryRepo) FindByID(_ context.Context, id uint) (*domain.Entry, error) {
	for _, e := range m.entries {
		if e.ID == id {
			return e, nil
		}
	}
	return nil, domain.ErrNotFound
}
func (m *mockEntryRepo) Create(_ context.Context, e *domain.Entry) error {
	m.nextID++
	e.ID = m.nextID
	m.entries = append(m.entries, e)
	return nil
}
func (m *mockEntryRepo) Update(_ context.Context, e *domain.Entry) error {
	for i, existing := range m.entries {
		if existing.ID == e.ID {
			m.entries[i] = e
			return nil
		}
	}
	return domain.ErrNotFound
}
func (m *mockEntryRepo) Delete(_ context.Context, id uint) error {
	for i, e := range m.entries {
		if e.ID == id {
			m.entries = append(m.entries[:i], m.entries[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}
func (m *mockEntryRepo) CountByModel(_ context.Context, modelID uint) (int64, error) {
	var count int64
	for _, e := range m.entries {
		if e.ContentModelID == modelID {
			count++
		}
	}
	return count, nil
}

// --- Tests ---

func blogModel() *domain.ContentModel {
	return &domain.ContentModel{
		ID:   1,
		Slug: "blog",
		Name: "Blog Post",
		Fields: []domain.FieldDefinition{
			{Name: "title", Type: domain.FieldTypeText, Required: true},
			{Name: "views", Type: domain.FieldTypeNumber, Required: false},
			{Name: "published_at", Type: domain.FieldTypeDate, Required: false},
			{Name: "body", Type: domain.FieldTypeMarkdown, Required: false},
		},
	}
}

func newSvc() *service.EntryService {
	cmRepo := &mockContentModelRepo{models: map[string]*domain.ContentModel{
		"blog": blogModel(),
	}}
	entryRepo := &mockEntryRepo{}
	return service.NewEntryService(cmRepo, entryRepo)
}

func TestEntryService_Create_Valid(t *testing.T) {
	svc := newSvc()
	content := domain.ContentData{
		"title": "Hello World",
		"views": float64(42),
	}
	entry, err := svc.Create(context.Background(), "blog", content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestEntryService_Create_MissingRequired(t *testing.T) {
	svc := newSvc()
	_, err := svc.Create(context.Background(), "blog", domain.ContentData{})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestEntryService_Create_WrongNumberType(t *testing.T) {
	svc := newSvc()
	_, err := svc.Create(context.Background(), "blog", domain.ContentData{
		"title": "Hello",
		"views": "not-a-number",
	})
	if !errors.Is(err, domain.ErrInvalidField) {
		t.Fatalf("expected ErrInvalidField, got %v", err)
	}
}

func TestEntryService_Create_ValidDate(t *testing.T) {
	svc := newSvc()
	_, err := svc.Create(context.Background(), "blog", domain.ContentData{
		"title":        "Hello",
		"published_at": "2024-01-15T10:00:00Z",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEntryService_Create_InvalidDate(t *testing.T) {
	svc := newSvc()
	_, err := svc.Create(context.Background(), "blog", domain.ContentData{
		"title":        "Hello",
		"published_at": "not-a-date",
	})
	if !errors.Is(err, domain.ErrInvalidField) {
		t.Fatalf("expected ErrInvalidField, got %v", err)
	}
}

func TestEntryService_Create_ModelNotFound(t *testing.T) {
	svc := newSvc()
	_, err := svc.Create(context.Background(), "nonexistent", domain.ContentData{})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestEntryService_Delete_NotFound(t *testing.T) {
	svc := newSvc()
	err := svc.Delete(context.Background(), 999)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
