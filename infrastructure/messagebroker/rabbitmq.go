package messagebroker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	conn      *amqp.Connection
	url       string
	queueName string
}

func NewRabbitMQPublisher(url, queueName string) (MessageBrokerInterface, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	return &RabbitMQPublisher{
		conn:      conn,
		url:       url,
		queueName: queueName,
	}, nil
}

func (publisher *RabbitMQPublisher) InitializeMessageBroker() {
	publisher.createQueue()
}

func (publisher *RabbitMQPublisher) createQueue() error {
	ch, err := publisher.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	// Declare the queue
	_, err = ch.QueueDeclare(
		publisher.queueName, // name
		true,                // durable
		false,               // auto-delete
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	return nil
}

func (publisher *RabbitMQPublisher) Publish(message interface{}) error {
	ch, err := publisher.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Publish the message
	err = ch.PublishWithContext(ctx,
		"",                  // exchange
		publisher.queueName, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("[x] Sent to queue %s: %s\n", publisher.queueName, body)
	return nil
}

func (publisher *RabbitMQPublisher) Close() {
	if publisher.conn != nil {
		publisher.conn.Close()
	}
}
