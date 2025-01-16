package messagebroker

import amqp "github.com/rabbitmq/amqp091-go"

type MessageBrokerConsumerInterface interface {
	Consume() (<-chan amqp.Delivery, error)
	InitializeMessageBroker()
	Close()
}
