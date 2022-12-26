package message

import (
	"os"
	"strconv"
	"sync"
	"time"
)

type MessageMetadata struct {
	BufferSize int
	Timeout    time.Duration
	WorkerSize int
}

type MessageConfig struct {
	metadata MessageMetadata
	once     sync.Once
}

// Private
func (messageConfig *MessageConfig) lazyInit() {
	messageConfig.once.Do(func() {
		bufferSize, err := strconv.Atoi(os.Getenv("BUFFER_SIZE"))
		if err != nil {
			panic(err)
		}
		numberofTimeoutSeconds, err := strconv.Atoi(os.Getenv("MESSAGE_TIMEOUT"))
		if err != nil {
			panic(err)
		}
		timeout := time.Duration(numberofTimeoutSeconds) * time.Second
		workerSize, err := strconv.Atoi(os.Getenv("WORKER_SIZE"))
		if err != nil {
			panic(err)
		}

		messageConfig.metadata.BufferSize = bufferSize
		messageConfig.metadata.Timeout = timeout
		messageConfig.metadata.WorkerSize = workerSize
	})
}

// Public
func (messageConfig *MessageConfig) GetMetadata() MessageMetadata {
	messageConfig.lazyInit()
	return messageConfig.metadata
}

var Config = &MessageConfig{}
