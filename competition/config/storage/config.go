package storage

import (
	"os"
	"strconv"
)

type StorageConfig struct {
	FileUploadTimeout int
	PhotoDir          string
	SubmissionDir     string
}

var currentStorageConfig *StorageConfig = nil

func Init() *StorageConfig {
	fileUploadTimeout, err := strconv.Atoi(os.Getenv("FILE_UPLOAD_TIMEOUT"))
	if err != nil {
		panic(err)
	}
	photoDir := os.Getenv("PHOTO_DIR")
	submissionDir := os.Getenv("SUBMISSION_DIR")

	return &StorageConfig{
		FileUploadTimeout: fileUploadTimeout,
		PhotoDir:          photoDir,
		SubmissionDir:     submissionDir,
	}
}

func GetStorageConfig() *StorageConfig {
	if currentStorageConfig == nil {
		currentStorageConfig = Init()
	}
	return currentStorageConfig
}
