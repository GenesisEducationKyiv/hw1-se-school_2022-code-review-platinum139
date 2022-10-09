package main

import (
	"github.com/fatih/color"
	"logs-printer/config"
	"logs-printer/pkg/message_brokers"
	"sync"
)

const logsExchangeName = "logs"

func main() {
	appConfig, err := config.NewAppConfig(".env")
	if err != nil {
		panic(err)
	}

	rabbitMqClient, err := message_brokers.NewRabbitMqClient(message_brokers.RabbitMqConfig{
		Host:     appConfig.RabbitMqHost,
		Port:     appConfig.RabbitMqPort,
		User:     appConfig.RabbitMqUserName,
		Password: appConfig.RabbitMqPassword,
		Timeout:  appConfig.RabbitMqTimeout,
		Exchange: logsExchangeName,
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := rabbitMqClient.CloseConnection(); err != nil {
			panic(err)
		}
	}()

	err = rabbitMqClient.CreateQueue(appConfig.LogLevel)
	if err != nil {
		panic(err)
	}

	logs, err := rabbitMqClient.GetFromQueue(appConfig.LogLevel)
	if err != nil {
		panic(err)
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()
		for log := range logs {
			printColor := color.New(color.FgYellow)
			printColor.Printf("%s", log.Timestamp)

			printColor = color.New(color.FgCyan)
			printColor.Printf(" [%s] ", appConfig.LogLevel)

			printColor = color.New(color.FgWhite)
			printColor.Printf("%s\n", string(log.Body))
		}
	}()

	waitGroup.Wait()
}
