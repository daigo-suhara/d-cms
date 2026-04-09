package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/daigo-suhara/d-cms/internal/domain"
	"github.com/daigo-suhara/d-cms/internal/repository"
	"github.com/google/uuid"
)

// StorageClient defines the interface for object storage operations.
type StorageClient interface {
	Upload(ctx context.Context, key string, r io.Reader, contentType string) (url string, err error)
	Delete(ctx context.Context, key string) error
}

type MediaService struct {
	mediaRepo repository.MediaRepository
	storage   StorageClient
}

func NewMediaService(mediaRepo repository.MediaRepository, storage StorageClient) *MediaService {
	return &MediaService{mediaRepo: mediaRepo, storage: storage}
}

func (s *MediaService) List(ctx context.Context) ([]domain.Media, error) {
	return s.mediaRepo.FindAll(ctx)
}

func (s *MediaService) Upload(ctx context.Context, fh *multipart.FileHeader) (*domain.Media, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("open uploaded file: %w", err)
	}
	defer f.Close()

	ext := filepath.Ext(fh.Filename)
	key := uuid.New().String() + ext
	contentType := fh.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	url, err := s.storage.Upload(ctx, key, f, contentType)
	if err != nil {
		return nil, fmt.Errorf("upload to storage: %w", err)
	}

	m := &domain.Media{
		Filename: fh.Filename,
		MimeType: contentType,
		Size:     fh.Size,
		URL:      url,
		Key:      key,
	}
	if err := s.mediaRepo.Create(ctx, m); err != nil {
		// Best-effort cleanup if DB save fails
		_ = s.storage.Delete(ctx, key)
		return nil, fmt.Errorf("save media metadata: %w", err)
	}
	return m, nil
}

func (s *MediaService) Delete(ctx context.Context, id uint) error {
	m, err := s.mediaRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.storage.Delete(ctx, m.Key); err != nil {
		return fmt.Errorf("delete from storage: %w", err)
	}
	return s.mediaRepo.Delete(ctx, id)
}
