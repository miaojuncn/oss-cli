package aliyun

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	bucket   = os.Getenv("bucket")
	uploader *AliOssStore
)

func TestUpload(t *testing.T) {
	// 获取一个断言实例
	should := assert.New(t)

	should.Equal(bucket, "appchain-production")

	err := uploader.Upload(bucket, "store.go", "store.go")
	if should.NoError(err) {
		// 	没有错误, 开启下一个步骤
		t.Log("upload ok")
	}
}

func TestUploadError(t *testing.T) {
	should := assert.New(t)

	err := uploader.Upload(bucket, "store.go", "store_xxx.go")
	should.Error(err, "open store_xxx.go: The system cannot find the file specified.")
}

func init() {
	ossStore, err := NewDefaultAliOssStore()
	if err != nil {
		panic(err)
	}
	uploader = ossStore
}
