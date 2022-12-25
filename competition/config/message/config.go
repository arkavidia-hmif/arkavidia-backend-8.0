package message

import (
	"os"
	"strconv"
	"time"
)

type MessageConfig struct {
	BufferSize int
	Timeout    time.Duration
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

	return &MessageConfig{
		BufferSize: bufferSize,
		Timeout:    timeout,
	}
}

func GetMessageConfig() *MessageConfig {
	if currentMessageConfig == nil {
		currentMessageConfig = Init()
	}
	return currentMessageConfig
}
