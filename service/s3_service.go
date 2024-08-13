package service

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type R2Service struct {
	s3Client   *s3.S3
	bucketName string
}

func NewR2Service() (*R2Service, error) {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("auto"),
        Endpoint: aws.String(os.Getenv("CLOUDFLARE_R2_ENDPOINT")),
        Credentials: credentials.NewStaticCredentials(
            os.Getenv("CLOUDFLARE_R2_ACCESS_KEY_ID"), 
            os.Getenv("CLOUDFLARE_R2_SECRET_ACCESS_KEY"), 
            "",
        ),
        S3ForcePathStyle: aws.Bool(true),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create session: %v", err)
    }

    svc := s3.New(sess)

    return &R2Service{
        s3Client:   svc,
        bucketName: os.Getenv("CLOUDFLARE_R2_BUCKET_NAME"),
    }, nil
}

// UploadFile uploads a file to Cloudflare R2
func (r2 *R2Service) UploadFile(key string, content []byte) (string, error) {
	_, err := r2.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(r2.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(content),
		ACL:    aws.String("public-read"), 
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to Cloudflare R2: %v", err)
	}

	return fmt.Sprintf("https://pub-736ef3be77f045e8ba550ae958fe7e1b.r2.dev/%s", key), nil
}
