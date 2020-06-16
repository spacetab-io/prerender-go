package s3

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	cfg "github.com/spacetab-io/prerender-go/configuration"
	"github.com/spacetab-io/prerender-go/pkg/models"
)

type storage struct {
	u   *s3manager.Uploader
	cfg cfg.S3Config
}

func NewStorage(cfg cfg.S3Config) storage { //nolint:golint
	s := new(storage)
	sess := session.Must(
		session.
			NewSession(aws.NewConfig().
				WithMaxRetries(3). //nolint:gomnd
				WithRegion(cfg.Region).
				WithCredentials(credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretKey, ""))))

	// Create an uploader with the session and default options
	s.u = s3manager.NewUploader(sess)
	s.cfg = cfg

	return *s
}

func (s storage) SaveData(pd *models.PageData) error {
	// Upload the file to S3.
	_, err := s.u.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s.cfg.Bucket),
		Key:         aws.String(s.cfg.BucketFolder + pd.FileName),
		Body:        bytes.NewReader(pd.Body),
		ContentType: aws.String("text/html; charset=utf-8"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	return nil
}
