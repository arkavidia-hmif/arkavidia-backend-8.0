package message

import (
	"os"
	"strconv"
	"time"
)

type MessageConfig struct {
	BufferSize int
	ReloadTime time.Duration
}

var currentMessageConfig *MessageConfig = nil

func Init() *MessageConfig {
	bufferSize, err := strconv.Atoi(os.Getenv("BUFFER_SIZE"))
	if err != nil {
		panic(err)
	}
	numberOfSeconds, err := strconv.Atoi(os.Getenv("RELOAD_TIME"))
	if err != nil {
		panic(err)
	}
	reloadTime := time.Duration(numberOfSeconds) * time.Second

	return &MessageConfig{
		BufferSize: bufferSize,
		ReloadTime: reloadTime,
	}
}

func GetMessageConfig() *MessageConfig {
	if currentMessageConfig == nil {
		currentMessageConfig = Init()
	}
	return currentMessageConfig
}
