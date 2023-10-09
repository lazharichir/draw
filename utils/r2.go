package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2 struct {
	s3 *s3.Client
}

func MustNewS3Client(
	accountID string,
	accessKeyID string,
	accessKeySecret string,
) *s3.Client {
	s3client, err := NewS3Client(accountID, accessKeyID, accessKeySecret)
	if err != nil {
		panic(err)
	}
	return s3client
}

func NewS3Client(
	accountID string,
	accessKeyID string,
	accessKeySecret string,
) (*s3.Client, error) {

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(`[R2] Connected to R2 Cloudflare Storage`, accountID, cfg.Region)
	return s3.NewFromConfig(cfg), nil
}

func New(
	accountID string,
	accessKeyID string,
	accessKeySecret string,
) (*R2, error) {
	s3client, err := NewS3Client(accountID, accessKeyID, accessKeySecret)
	return &R2{
		s3: s3client,
	}, err
}

func (r2 *R2) DeleteObject(ctx context.Context, bucket string, key string) error {
	output, err := r2.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	_ = output

	if err != nil {
		return fmt.Errorf("delete object failed: %w", err)
	}

	return nil
}

func (r2 *R2) PutBytes(ctx context.Context, bucket string, key string, data []byte) error {
	log.Printf("[R2] PUT %s::%s \n", bucket, key)

	body := bytes.NewBuffer(data)

	output, err := r2.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   body,
	})

	_ = output

	if err != nil {
		return fmt.Errorf("put bytes failed: %w", err)
	}

	return nil
}

func (r2 *R2) GetBytes(ctx context.Context, bucket string, key string) ([]byte, error) {
	log.Printf("[R2] GET %s::%s \n", bucket, key)
	output, err := r2.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("get bytes failed: %w", err)
	}

	data, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, fmt.Errorf("read bytes failed: %w", err)
	}

	return data, nil
}
