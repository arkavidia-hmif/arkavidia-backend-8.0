package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"

	storageConfig "arkavidia-backend-8.0/competition/config/storage"
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

func UploadFile(client *storage.Client, filename string, uploadPath string, file multipart.File) error {
	config := storageConfig.GetStorageConfig()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Duration(config.FileTimeout).Seconds()))
	defer cancel()

	storageWriter := client.Bucket(config.BucketName).Object(fmt.Sprintf("%s/%s", uploadPath, filename)).NewWriter(ctx)
	if _, err := io.Copy(storageWriter, file); err != nil {
		return err
	}
	if err := storageWriter.Close(); err != nil {
		return err
	}

	return nil
}

func DownloadFile(client *storage.Client, filename string, downloadPath string) (io.Writer, error) {
	config := storageConfig.GetStorageConfig()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Duration(config.FileTimeout).Seconds()))
	defer cancel()

	storageReader, err := client.Bucket(config.BucketName).Object(fmt.Sprintf("%s/%s", downloadPath, filename)).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer storageReader.Close()

	var file io.Writer
	if _, err := io.Copy(file, storageReader); err != nil {
		return nil, err
	}

	return file, nil
}

func DeleteFile(client *storage.Client, filename string, deletePath string) error {
	config := storageConfig.GetStorageConfig()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Duration(config.FileTimeout).Seconds()))
	defer cancel()

	object := client.Bucket(config.BucketName).Object(fmt.Sprintf("%s/%s", deletePath, filename))
	// Handles Race Condition
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return err
	}
	object = object.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := object.Delete(ctx); err != nil {
		return err
	}

	return nil
}
