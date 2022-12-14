package storage

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/sync/singleflight"

	storageConfig "arkavidia-backend-8.0/competition/config/storage"
)

// TODO: Tambahkan duplicate function call suppression mechanism menggunakan lib Singleflight
// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-singleflight.html
// ASSIGNED TO: @akbarmridho
// STATUS: DONE

type StorageClient struct {
	client       *storage.Client
	requestGroup singleflight.Group
	once         sync.Once
}

func (storageClient *StorageClient) lazyInit() {
	storageClient.once.Do(func() {
		client, err := storage.NewClient(context.Background())
		if err != nil {
			panic(err)
		}

		storageClient.client = client
	})
}

// Public
func (storageClient *StorageClient) DownloadFile(filename string, downloadPath string) (io.Writer, error) {
	storageClient.lazyInit()

	// Duplicate Function Call Suppression Mechanism
	v, err, _ := storageClient.requestGroup.Do(filename, func() (interface{}, error) {
		config := storageConfig.Config.GetMetadata()
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.FileTimeout)*time.Second)
		defer cancel()

		storageReader, err := storageClient.client.Bucket(config.BucketName).Object(fmt.Sprintf("%s/%s", downloadPath, filename)).NewReader(ctx)
		if err != nil {
			return nil, err
		}
		defer storageReader.Close()

		var file io.Writer
		if _, err := io.Copy(file, storageReader); err != nil {
			return nil, err
		}

		return file, nil
	})
	if err != nil {
		return nil, err
	}

	file := v.(io.Writer)

	return file, nil
}

func (storageClient *StorageClient) UploadFile(filename string, uploadPath string, content io.Reader) error {
	storageClient.lazyInit()

	config := storageConfig.Config.GetMetadata()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.FileTimeout)*time.Second)
	defer cancel()

	storageWriter := storageClient.client.Bucket(config.BucketName).Object(fmt.Sprintf("%s/%s", uploadPath, filename)).NewWriter(ctx)
	if _, err := io.Copy(storageWriter, content); err != nil {
		return err
	}
	if err := storageWriter.Close(); err != nil {
		return err
	}

	return nil
}

// Currently Unused (Soft Deletion)
func (storageClient *StorageClient) DeleteFile(filename string, deletePath string) error {
	storageClient.lazyInit()

	config := storageConfig.Config.GetMetadata()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.FileTimeout)*time.Second)
	defer cancel()

	object := storageClient.client.Bucket(config.BucketName).Object(fmt.Sprintf("%s/%s", deletePath, filename))
	if err := object.Delete(ctx); err != nil {
		return err
	}

	return nil
}

var Client = &StorageClient{}
