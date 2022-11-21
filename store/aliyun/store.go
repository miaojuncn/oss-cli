package aliyun

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliOssStore struct {
	client *oss.Client
}

func (s *AliOssStore) Upload(bucketName string, objectKey string, fileName string) error {
	bucket, err := s.client.Bucket(bucketName)
	if err != nil {
		return err
	}

	if err := bucket.PutObjectFromFile(objectKey, fileName); err != nil {
		return err
	}

	// 打印下载链接
	downloadURL, err := bucket.SignURL(objectKey, oss.HTTPGet, 60*60*24)
	if err != nil {
		return err
	}
	fmt.Printf("文件下载URL: %s \n", downloadURL)
	fmt.Println("请在1天之内下载.")
	return nil
}

func NewAliOssStore(endpoint, accessKey, accessSecret string) (*AliOssStore, error) {
	client, err := oss.New(endpoint, accessKey, accessSecret)
	if err != nil {
		return nil, err
	}
	return &AliOssStore{client: client}, nil
}
