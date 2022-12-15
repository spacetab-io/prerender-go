package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spacetab-io/prerender-go/configuration"
	"github.com/spacetab-io/prerender-go/pkg/models"
)

type storage struct {
	client *s3.Client
	cfg    configuration.S3Config
}

var ErrUnknownEndpointRequest = errors.New("unknown endpoint requested")

//nolint:revive // we need it here
func NewStorage(cfg configuration.S3Config) (*storage, error) {
	// Создаем кастомный обработчик эндпоинтов, который для сервиса S3 и региона ru-central1 выдаст корректный URL
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				PartitionID:       cfg.Endpoint.PartitionID,
				URL:               strings.TrimRight(cfg.Endpoint.URL, "/"),
				SigningName:       cfg.Endpoint.SigningName,
				SigningRegion:     cfg.Endpoint.SigningRegion,
				SigningMethod:     cfg.Endpoint.SigningMethod,
				HostnameImmutable: true,
			}, nil
		}

		return aws.Endpoint{}, ErrUnknownEndpointRequest
	})

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.Access.AccessKeyID,
			cfg.Access.SecretKey,
			cfg.Access.Token,
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("config load error: %w", err)
	}

	return &storage{client: s3.NewFromConfig(awsCfg), cfg: cfg}, nil
}

func (s *storage) SaveData(ctx context.Context, pd *models.PageData) error {
	// Upload the file to S3.
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(strings.Trim(s.cfg.Bucket.Name, "/")),
		Key:         aws.String(fmt.Sprintf("%s/%s", strings.Trim(s.cfg.Bucket.Folder, "/"), pd.FileName)),
		Body:        bytes.NewReader(pd.Body),
		ContentType: aws.String("text/html; charset=utf-8"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %w", err)
	}

	return nil
}
