package store

type Uploader interface {
	Upload(bucketName string, objectKey string, fileName string) error
}
