package s3

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	cfg "github.com/spacetab-io/prerender-go/configuration"
	"github.com/spacetab-io/prerender-go/pkg/models"
)

type storage struct {
	client *s3.Client
	cfg    cfg.S3Config
}

func NewStorage(cfg cfg.S3Config) (*storage, error) { //nolint:golint
	s := new(storage)
	// Создаем кастомный обработчик эндпоинтов, который для сервиса S3 и региона ru-central1 выдаст корректный URL
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				PartitionID:       cfg.Endpoint.PartitionID,
				URL:               cfg.Endpoint.URL,
				SigningName:       cfg.Endpoint.SigningName,
				SigningRegion:     cfg.Endpoint.SigningRegion,
				SigningMethod:     cfg.Endpoint.SigningMethod,
				HostnameImmutable: true,
			}, nil
		}

		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	// Подгружаем конфигрурацию из ~/.aws/*
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

	s.client = s3.NewFromConfig(awsCfg)
	s.cfg = cfg

	return s, nil
}

func (s storage) SaveData(pd *models.PageData) error {
	// Upload the file to S3.
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.cfg.Bucket.Name),
		Key:         aws.String(s.cfg.Bucket.Folder + pd.FileName),
		Body:        bytes.NewReader(pd.Body),
		ContentType: aws.String("text/html; charset=utf-8"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	return nil
}
