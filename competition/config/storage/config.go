package storage

import (
	"os"
	"strconv"
	"sync"
)

type StorageMetadata struct {
	FileTimeout   int
	StorageHost   string
	BucketName    string
	PhotoDir      string
	SubmissionDir string
}

type StorageConfig struct {
	metadata StorageMetadata
	once     sync.Once
}

// Private
func (storageConfig *StorageConfig) lazyInit() {
	storageConfig.once.Do(func() {
		fileTimeout, err := strconv.Atoi(os.Getenv("FILE_TIMEOUT"))
		if err != nil {
			panic(err)
		}
		storageHost := os.Getenv("STORAGE_HOST")
		bucketName := os.Getenv("BUCKET_NAME")
		photoDir := os.Getenv("PHOTO_DIR")
		submissionDir := os.Getenv("SUBMISSION_DIR")

		storageConfig.metadata.FileTimeout = fileTimeout
		storageConfig.metadata.StorageHost = storageHost
		storageConfig.metadata.BucketName = bucketName
		storageConfig.metadata.PhotoDir = photoDir
		storageConfig.metadata.SubmissionDir = submissionDir
	})
}

// Public
func (storageConfig *StorageConfig) GetMetadata() StorageMetadata {
	storageConfig.lazyInit()
	return storageConfig.metadata
}

var Config = &StorageConfig{}
