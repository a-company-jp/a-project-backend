package gcs

import (
	"a-project-backend/pkg/config"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/option"
)

type GCS struct {
	bucketName string
	client     *storage.Client
}

func NewGCS() (*GCS, error) {
	conf := config.Get()
	bucketName := conf.Application.GCS.BucketName

	ctx := context.Background()
	var client *storage.Client
	var err error
	if conf.Infrastructure.GoogleCloud.UseCredentialsFile {
		client, err = storage.NewClient(ctx,
			option.WithCredentialsFile(conf.Infrastructure.GoogleCloud.CredentialsFilePath))
		if err != nil {
			return nil, fmt.Errorf("failed to create GCS client with credentials file: %w", err)
		}
	} else {
		client, err = storage.NewClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create GCS client: %w", err)
		}
	}
	return &GCS{
		bucketName: bucketName,
		client:     client,
	}, nil
}

func (g *GCS) Upload(ctx context.Context, objectName string, data []byte) error {
	wc := g.client.Bucket(g.bucketName).Object(objectName).NewWriter(ctx)

	if _, err := wc.Write(data); err != nil {
		return fmt.Errorf("failed to write to GCS: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return nil
}
