package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/daigo-suhara/d-cms/config"
)

type R2Client struct {
	client     *s3.Client
	bucket     string
	publicBase string
}

func NewR2Client(cfg *config.Config) (*R2Client, error) {
	if cfg.R2Endpoint == "" {
		return &R2Client{}, nil // No-op client when R2 is not configured
	}

	awsCfg := aws.Config{
		Region:      cfg.AWSRegion,
		Credentials: credentials.NewStaticCredentialsProvider(cfg.AWSKeyID, cfg.AWSKeySecret, ""),
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.R2Endpoint)
		o.UsePathStyle = true
	})

	return &R2Client{
		client:     client,
		bucket:     cfg.R2BucketName,
		publicBase: cfg.R2PublicBase,
	}, nil
}

func (r *R2Client) Upload(ctx context.Context, key string, body io.Reader, contentType string) (string, error) {
	if r.client == nil {
		return "", fmt.Errorf("R2 storage is not configured")
	}
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("R2 PutObject %q: %w", key, err)
	}
	return r.publicBase + "/" + key, nil
}

func (r *R2Client) Delete(ctx context.Context, key string) error {
	if r.client == nil {
		return fmt.Errorf("R2 storage is not configured")
	}
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("R2 DeleteObject %q: %w", key, err)
	}
	return nil
}

func (r *R2Client) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	if r.client == nil {
		return nil, fmt.Errorf("R2 storage is not configured")
	}
	result, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("R2 GetObject %q: %w", key, err)
	}
	return result.Body, nil
}

func (r *R2Client) ListObjects(ctx context.Context, prefix string) ([]string, error) {
	if r.client == nil {
		return nil, fmt.Errorf("R2 storage is not configured")
	}
	paginator := s3.NewListObjectsV2Paginator(r.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(r.bucket),
		Prefix: aws.String(prefix),
	})
	var keys []string
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("R2 ListObjectsV2 prefix=%q: %w", prefix, err)
		}
		for _, obj := range page.Contents {
			keys = append(keys, aws.ToString(obj.Key))
		}
	}
	return keys, nil
}
