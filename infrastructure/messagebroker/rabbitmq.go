package messagebroker

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
}

func NewRabbitMQConsumer(url, queueName string) (MessageBrokerConsumerInterface, error) {
	connection, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		connection.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &RabbitMQConsumer{
		connection: connection,
		channel:    channel,
		queueName:  queueName,
	}, nil
}

func (consumer *RabbitMQConsumer) InitializeMessageBroker() {
}

func (consumer *RabbitMQConsumer) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := consumer.channel.Consume(
		consumer.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start consuming messages: %w", err)
	}

	return msgs, nil
}

func (consumer *RabbitMQConsumer) Close() {
	if consumer.channel != nil {
		consumer.channel.Close()
	}
	if consumer.connection != nil {
		consumer.connection.Close()
	}
}
