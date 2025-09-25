package r2

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config" // "config" 충돌 방지 위해 별칭 사용
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pro200/go-utils"
)

type Config struct {
	Name            string // storage 이름 (기본값: "main")
	AccountId       string
	AccessKeyID     string
	SecretAccessKey string
}

type Storage struct {
	client *s3.Client
}

var (
	Storages = make(map[string]*Storage)
	dbMu     sync.RWMutex // 동시성 안전 보장
)

func New(config Config) (*Storage, error) {
	if config.Name == "" {
		config.Name = "main"
	}

	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.AccessKeyID, config.SecretAccessKey, "")),
		awsConfig.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", config.AccountId))
	})

	storage := &Storage{client: client}

	dbMu.Lock()
	Storages[config.Name] = storage
	dbMu.Unlock()

	return storage, nil
}

func GetStorage(name ...string) (*Storage, error) {
	dbMu.RLock()
	defer dbMu.RUnlock()

	if len(Storages) == 0 {
		return nil, errors.New("no storages available")
	}

	stName := "main"
	if len(name) > 0 {
		stName = name[0]
	}

	db, ok := Storages[stName]
	if !ok {
		return nil, fmt.Errorf("storage %s not found", stName)
	}
	return db, nil
}

func (s *Storage) Info(bucket, key string) (*s3.HeadObjectOutput, error) {
	return s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
}

func (s *Storage) List(bucket, prefix string, length int, token ...string) (list []string, nextToken string, err error) {
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

	output, err := s.client.ListObjectsV2(context.TODO(), &options)
	if err != nil {
		return list, nextToken, err
	}

	for _, obj := range output.Contents {
		list = append(list, aws.ToString(obj.Key))
	}

	nextToken = aws.ToString(output.NextContinuationToken)
	return list, nextToken, nil
}

func (s *Storage) Upload(bucket, path, key string, forceType ...string) error {
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

	uploader := manager.NewUploader(s.client)
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
	result, err := s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	// TODO: 업로드 실패한 파일을 삭제
	if len(file) != int(*result.ContentLength) {
		return errors.New("upload failed")
	}

	return nil
}

func (s *Storage) Delete(bucket, key string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return err
}

func (s *Storage) Download(bucket, key, targetPath string) error {
	// Set up the local file
	fd, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer fd.Close()

	downloader := manager.NewDownloader(s.client)
	_, err = downloader.Download(context.TODO(), fd,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	return err
}
