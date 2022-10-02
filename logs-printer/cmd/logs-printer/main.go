package main

import (
	"github.com/fatih/color"
	"logs-printer/config"
	"logs-printer/pkg/message_brokers"
	"sync"
)

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
		Exchange: "logs",
	})
	if err != nil {
		panic(err)
	}
	defer rabbitMqClient.CloseConnection()

	err = rabbitMqClient.CreateQueue(appConfig.LogLevel)
	if err != nil {
		panic(err)
	}

	logs, err := rabbitMqClient.GetFromQueue(appConfig.LogLevel)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for log := range logs {
			c := color.New(color.FgYellow)
			c.Printf("%s", log.Timestamp)

			c = color.New(color.FgCyan)
			c.Printf(" [%s] ", appConfig.LogLevel)

			c = color.New(color.FgWhite)
			c.Printf("%s\n", string(log.Body))
		}
	}()

	wg.Wait()
}
