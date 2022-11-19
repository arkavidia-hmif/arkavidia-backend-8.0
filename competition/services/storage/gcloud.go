package storage

import (
	"context"
	"mime/multipart"

	"cloud.google.com/go/storage"
)

var currentClient *storage.Client = nil

func Init() *storage.Client {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		panic(err)
	}

	return client
}

func GetClient() *storage.Client {
	if currentClient == nil {
		currentClient = Init()
	}
	return currentClient
}

func UploadFile(client *storage.Client, filename string, uploadPath string, file multipart.File) error

func DeleteFile(client *storage.Client, filename string, deletePath string) error
