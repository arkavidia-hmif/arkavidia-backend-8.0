package storage

import (
	"os"
	"strconv"
)

type StorageConfig struct {
	FileTimeout   int
	StorageHost   string
	BucketName    string
	PhotoDir      string
	SubmissionDir string
}

var currentStorageConfig *StorageConfig = nil

func Init() *StorageConfig {
	fileTimeout, err := strconv.Atoi(os.Getenv("FILE_TIMEOUT"))
	if err != nil {
		panic(err)
	}
	storageHost := os.Getenv("STORAGE_HOST")
	bucketName := os.Getenv("BUCKET_NAME")
	photoDir := os.Getenv("PHOTO_DIR")
	submissionDir := os.Getenv("SUBMISSION_DIR")

	return &StorageConfig{
		FileTimeout:   fileTimeout,
		StorageHost:   storageHost,
		BucketName:    bucketName,
		PhotoDir:      photoDir,
		SubmissionDir: submissionDir,
	}
}

func GetStorageConfig() *StorageConfig {
	if currentStorageConfig == nil {
		currentStorageConfig = Init()
	}
	return currentStorageConfig
}
