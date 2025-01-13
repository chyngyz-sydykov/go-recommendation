package messagebroker

type MessageBrokerInterface interface {
	Publish(message interface{}) error
	InitializeMessageBroker()
	Close()
}
