package message_brokers

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type RabbitMqConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Exchange string
}

type RabbitMqClient struct {
	config     RabbitMqConfig
	connection *amqp.Connection
	channel    *amqp.Channel
	queues     []amqp.Queue
}

func (c *RabbitMqClient) CreateQueue(name string) error {
	queue, err := c.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return err
	}

	err = c.channel.QueueBind(
		queue.Name, // queue name
		queue.Name, // routing key
		"logs",     // exchange
		false,      // noWait
		nil,        // args
	)
	if err != nil {
		return err
	}

	c.queues = append(c.queues, queue)
	return nil
}

func (c *RabbitMqClient) SendToQueue(queueName string, messageBody string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         []byte(messageBody),
	}

	err := c.channel.PublishWithContext(
		ctx,
		c.config.Exchange, // exchange
		queueName,         // routing key
		true,              // mandatory
		false,             // immediate
		message)

	return err
}

func (c *RabbitMqClient) GetFromQueue(queueName string) (<-chan amqp.Delivery, error) {
	delivery, err := c.channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	return delivery, err
}

func (c *RabbitMqClient) CloseConnection() error {
	if err := c.connection.Close(); err != nil {
		return err
	}
	if err := c.channel.Close(); err != nil {
		return err
	}
	return nil
}

func NewRabbitMqClient(config RabbitMqConfig) (*RabbitMqClient, error) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.User, config.Password, config.Host, config.Port)

	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = channel.ExchangeDeclare(
		config.Exchange, // name
		"direct",        // kind
		true,            //durable
		false,           //autoDelete
		false,           // internal
		false,           // noWait
		nil,             // args
	)
	if err != nil {
		return nil, err
	}

	client := &RabbitMqClient{
		config:     config,
		connection: conn,
		channel:    channel,
	}

	return client, err
}
