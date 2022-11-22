package aliyun

import (
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-playground/validator"
)

var (
	validate = validator.New()
)

type Options struct {
	Endpoint     string `validate:"required"`
	AccessKey    string `validate:"required"`
	AccessSecret string `validate:"required"`
}

func (o *Options) Validate() error {
	return validate.Struct(o)
}

type AliOssStore struct {
	client   *oss.Client
	listener *OssProgressListener
}

func NewDefaultAliOssStore() (*AliOssStore, error) {
	return NewAliOssStore(&Options{
		Endpoint:     os.Getenv("endpoint"),
		AccessKey:    os.Getenv("accessKey"),
		AccessSecret: os.Getenv("accessSecret"),
	})
}

func (s *AliOssStore) GetBucket(bucketName string) (*oss.Bucket, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("upload bucket name required")
	}
	bucket, err := s.client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func (s *AliOssStore) Upload(bucketName string, objectKey string, fileName string) (downloadUrl string, err error) {
	bucket, err := s.GetBucket(bucketName)
	if err != nil {
		return
	}

	if err = bucket.PutObjectFromFile(objectKey, fileName, oss.Progress(&OssProgressListener{})); err != nil {
		return
	}
	// 打印下载链接
	downloadUrl, err = bucket.SignURL(objectKey, oss.HTTPGet, 60*60*24)
	if err != nil {
		return
	}
	return downloadUrl, nil
}

func NewAliOssStore(opts *Options) (*AliOssStore, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	client, err := oss.New(opts.Endpoint, opts.AccessKey, opts.AccessSecret)
	if err != nil {
		return nil, err
	}
	return &AliOssStore{client: client, listener: NewOssProgressListener()}, nil
}
