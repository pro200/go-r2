package r2

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pro200/go-utils"
)

var s3Client *s3.Client

type Config struct {
	AccountId       string
	AccessKeyID     string
	SecretAccessKey string
}

func Init(option Config) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(option.AccessKeyID, option.SecretAccessKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", option.AccountId))
	})

	return nil
}

func Info(bucket, key string) (*s3.HeadObjectOutput, error) {
	return s3Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
}

func List(bucket, prefix string, length int, token ...string) (list []string, nextToken string, err error) {
	// up to 1,000 keys
	if length > 1000 {
		length = 1000
	}

	options := s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int32(int32(length)),
	}

	// ContinuationToken
	// A token to specify where to start paginating. This is the NextContinuationToken from a previously truncated response.
	if len(token) > 0 {
		options.ContinuationToken = aws.String(token[0])
	}

	output, err := s3Client.ListObjectsV2(context.TODO(), &options)

	if err != nil {
		return list, nextToken, err
	}

	for _, obj := range output.Contents {
		list = append(list, aws.ToString(obj.Key))
	}

	nextToken = aws.ToString(output.NextContinuationToken)
	return list, nextToken, nil
}

func Upload(bucket, path, key string, forceType ...string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if len(file) == 0 {
		return errors.New("zero size file")
	}

	contentType := utils.ContentType(path)
	if len(forceType) > 0 {
		contentType = forceType[0]
	}

	uploader := manager.NewUploader(s3Client)
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(file),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return err
	}

	// 업로드된 용량 비교
	result, err := s3Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil || len(file) != int(*result.ContentLength) {
		return errors.New("upload failed")
	}

	return err
}

func Delete(bucket, key string) error {
	_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return err
}

func Download(bucket, key, targetPath string) error {
	// Set up the local file
	fd, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer fd.Close()

	downloader := manager.NewDownloader(s3Client)
	_, err = downloader.Download(context.TODO(), fd,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	return err
}
