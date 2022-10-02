package common

import (
	"fmt"
	"subscribers-service/pkg/message_brokers"
)

const (
	debugLogQueue = "debug"
	infoLogQueue  = "info"
	errorLogQueue = "error"
)

type LoggerImpl struct {
	rabbitMqClient *message_brokers.RabbitMqClient
}

func (log *LoggerImpl) Debugf(format string, args ...interface{}) {
	err := log.rabbitMqClient.CreateQueue(debugLogQueue)
	if err != nil {
		fmt.Println("Failed to create debug queue")
	}

	message := fmt.Sprintf(format, args)
	err = log.rabbitMqClient.SendToQueue(debugLogQueue, message)
	if err != nil {
		fmt.Println("Failed to write debug log")
	}
}

func (log *LoggerImpl) Infof(format string, args ...interface{}) {
	err := log.rabbitMqClient.CreateQueue(infoLogQueue)
	if err != nil {
		fmt.Println("Failed to create info queue")
	}

	message := fmt.Sprintf(format, args)
	err = log.rabbitMqClient.SendToQueue(infoLogQueue, message)
	if err != nil {
		fmt.Println("Failed to write info log")
	}
}

func (log *LoggerImpl) Errorf(format string, args ...interface{}) {
	err := log.rabbitMqClient.CreateQueue(errorLogQueue)
	if err != nil {
		fmt.Println("Failed to create error queue")
	}

	message := fmt.Sprintf(format, args)
	err = log.rabbitMqClient.SendToQueue(errorLogQueue, message)
	if err != nil {
		fmt.Println("Failed to write error log")
	}
}

func NewLogger(client *message_brokers.RabbitMqClient) *LoggerImpl {
	return &LoggerImpl{
		rabbitMqClient: client,
	}
}
