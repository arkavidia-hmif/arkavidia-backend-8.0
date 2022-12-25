package message

import (
	"os"
	"strconv"
	"time"
)

type MessageConfig struct {
	BufferSize int
	Timeout    time.Duration
	WorkerSize int
}

var currentMessageConfig *MessageConfig = nil

func Init() *MessageConfig {
	bufferSize, err := strconv.Atoi(os.Getenv("BUFFER_SIZE"))
	if err != nil {
		panic(err)
	}
	numberofTimeoutSeconds, err := strconv.Atoi(os.Getenv("MESSAGE_TIMEOUT"))
	if err != nil {
		panic(err)
	}
	timeout := time.Duration(numberofTimeoutSeconds) * time.Second
	numberOfWorker, err := strconv.Atoi(os.Getenv("WORKER_SIZE"))
	if err != nil {
		panic(err)
	}

	return &MessageConfig{
		BufferSize: bufferSize,
		Timeout:    timeout,
		WorkerSize: numberOfWorker,
	}
}

func GetMessageConfig() *MessageConfig {
	if currentMessageConfig == nil {
		currentMessageConfig = Init()
	}
	return currentMessageConfig
}
