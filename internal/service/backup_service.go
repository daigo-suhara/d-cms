package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/daigo-suhara/d-cms/internal/domain"
	"github.com/daigo-suhara/d-cms/internal/repository"
)

// BackupStorage extends StorageClient with download and list capabilities.
type BackupStorage interface {
	Upload(ctx context.Context, key string, r io.Reader, contentType string) (url string, err error)
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	ListObjects(ctx context.Context, prefix string) ([]string, error)
}

const backupPrefix = "backups/"

// BackupEntry stores entry data with the content model slug for portability.
type BackupEntry struct {
	ContentModelSlug string            `json:"content_model_slug"`
	Content          domain.ContentData `json:"content"`
}

// BackupData is the top-level structure serialized to JSON in R2.
type BackupData struct {
	Version       string               `json:"version"`
	CreatedAt     time.Time            `json:"created_at"`
	ContentModels []domain.ContentModel `json:"content_models"`
	Entries       []BackupEntry        `json:"entries"`
	Media         []domain.Media       `json:"media"`
}

// BackupInfo describes a single backup file stored in R2.
type BackupInfo struct {
	Key       string
	Name      string
	CreatedAt time.Time
}

type BackupService struct {
	cmRepo    repository.ContentModelRepository
	entryRepo repository.EntryRepository
	mediaRepo repository.MediaRepository
	storage   BackupStorage
}

func NewBackupService(
	cmRepo repository.ContentModelRepository,
	entryRepo repository.EntryRepository,
	mediaRepo repository.MediaRepository,
	storage BackupStorage,
) *BackupService {
	return &BackupService{
		cmRepo:    cmRepo,
		entryRepo: entryRepo,
		mediaRepo: mediaRepo,
		storage:   storage,
	}
}

// CreateBackup dumps all DB data to a JSON file and uploads it to R2.
// Returns the storage key of the created backup.
func (s *BackupService) CreateBackup(ctx context.Context) (string, error) {
	models, err := s.cmRepo.FindAll(ctx)
	if err != nil {
		return "", fmt.Errorf("fetch content models: %w", err)
	}

	var entries []BackupEntry
	for _, m := range models {
		es, err := s.entryRepo.FindAllByModel(ctx, m.ID)
		if err != nil {
			return "", fmt.Errorf("fetch entries for model %q: %w", m.Slug, err)
		}
		for _, e := range es {
			entries = append(entries, BackupEntry{
				ContentModelSlug: m.Slug,
				Content:          e.Content,
			})
		}
	}

	media, err := s.mediaRepo.FindAll(ctx)
	if err != nil {
		return "", fmt.Errorf("fetch media: %w", err)
	}

	data := BackupData{
		Version:       "1",
		CreatedAt:     time.Now().UTC(),
		ContentModels: models,
		Entries:       entries,
		Media:         media,
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal backup data: %w", err)
	}

	key := backupPrefix + "backup-" + data.CreatedAt.Format("2006-01-02T15-04-05Z") + ".json"
	if _, err := s.storage.Upload(ctx, key, bytes.NewReader(b), "application/json"); err != nil {
		return "", fmt.Errorf("upload backup: %w", err)
	}
	return key, nil
}

// ListBackups returns metadata for all backup files stored in R2.
func (s *BackupService) ListBackups(ctx context.Context) ([]BackupInfo, error) {
	keys, err := s.storage.ListObjects(ctx, backupPrefix)
	if err != nil {
		return nil, fmt.Errorf("list backups: %w", err)
	}

	infos := make([]BackupInfo, 0, len(keys))
	for _, key := range keys {
		name := strings.TrimPrefix(key, backupPrefix)
		// Parse creation time from filename: backup-2006-01-02T15-04-05Z.json
		var t time.Time
		trimmed := strings.TrimPrefix(name, "backup-")
		trimmed = strings.TrimSuffix(trimmed, ".json")
		t, _ = time.Parse("2006-01-02T15-04-05Z", trimmed)
		infos = append(infos, BackupInfo{Key: key, Name: name, CreatedAt: t})
	}
	return infos, nil
}

// RestoreBackup downloads a backup from R2 and imports its data into the DB.
// ContentModels are skipped if their slug already exists.
// Entries and media are always created.
func (s *BackupService) RestoreBackup(ctx context.Context, key string) error {
	rc, err := s.storage.Download(ctx, key)
	if err != nil {
		return fmt.Errorf("download backup %q: %w", key, err)
	}
	defer func() {
		_ = rc.Close()
	}()

	var data BackupData
	if err := json.NewDecoder(rc).Decode(&data); err != nil {
		return fmt.Errorf("decode backup: %w", err)
	}

	// Restore content models, building old-slug → new-ID map.
	slugToID := make(map[string]uint, len(data.ContentModels))
	for i := range data.ContentModels {
		m := &data.ContentModels[i]
		existing, err := s.cmRepo.FindBySlug(ctx, m.Slug)
		if err == nil {
			// Already exists — reuse its ID.
			slugToID[m.Slug] = existing.ID
			continue
		}
		m.ID = 0 // Let DB assign a new ID.
		if err := s.cmRepo.Create(ctx, m); err != nil {
			return fmt.Errorf("create content model %q: %w", m.Slug, err)
		}
		slugToID[m.Slug] = m.ID
	}

	// Restore entries.
	for _, be := range data.Entries {
		modelID, ok := slugToID[be.ContentModelSlug]
		if !ok {
			continue // Skip entries whose model was not restored.
		}
		e := &domain.Entry{
			ContentModelID: modelID,
			Content:        be.Content,
		}
		if err := s.entryRepo.Create(ctx, e); err != nil {
			return fmt.Errorf("create entry for model %q: %w", be.ContentModelSlug, err)
		}
	}

	// Restore media metadata (skip if URL already exists).
	existing, err := s.mediaRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("fetch existing media: %w", err)
	}
	existingURLs := make(map[string]struct{}, len(existing))
	for _, m := range existing {
		existingURLs[m.URL] = struct{}{}
	}
	for i := range data.Media {
		m := &data.Media[i]
		if _, found := existingURLs[m.URL]; found {
			continue
		}
		m.ID = 0
		if err := s.mediaRepo.Create(ctx, m); err != nil {
			return fmt.Errorf("create media %q: %w", m.Filename, err)
		}
	}

	return nil
}
