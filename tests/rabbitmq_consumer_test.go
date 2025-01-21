package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chyngyz-sydykov/go-recommendation/application"
	"github.com/chyngyz-sydykov/go-recommendation/application/handlers"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/config"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/db/models"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/logger"
	"github.com/chyngyz-sydykov/go-recommendation/internal/recommendation"
	"github.com/stretchr/testify/mock"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (suite *IntegrationSuite) TestRecommendationUpsert_ShouldSuccessfullySaveWithRealRabbitMQ() {
	testCases := []struct {
		name            string
		BookId          int
		Event           string
		ResultingPoints int
	}{

		{"InitialBookUpdate", 3, "bookUpdated", 1},
		{"RecurringBookUpdateOnUpdatedBook", 3, "bookUpdated", 2},
		{"InitialBookRated", 4, "bookRated", 3},
		{"RecurringBookUpdateOnRatedBook", 4, "bookUpdated", 4},
		{"RecurringBookRateOnRatedBook", 4, "bookRated", 7},
	}
	sqlLiteDB := initializeSqlLite()

	logger := logger.NewLogger()
	commonHandler := handlers.NewCommonHandler(logger)

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not config: %v", err)
	}
	consumer := application.InitializeRabbitMqConsumer(config, logger)

	rabbitMQURL := "amqp://" + config.RabbitMqUser + ":" + config.RabbitMqPassword + "@" + config.RabbitMqContainerName + ":5672/"
	con, ch := newPublisher(rabbitMQURL, config.RabbitMqQueueName)

	defer closePublisher(con, ch)

	for _, testCase := range testCases {
		suite.T().Run(testCase.name, func(t *testing.T) {
			bookMessage := recommendation.BookMessage{
				BookId: testCase.BookId,
				Event:  testCase.Event,
			}

			err := publish(ch, config.RabbitMqQueueName, bookMessage)
			if err != nil {
				log.Fatalf("failed to publish message: %v", err)
			}
		})
	}
	app := provideDependencies(suite, consumer, commonHandler)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		fmt.Println("Waiting for 2 seconds")
		time.Sleep(2 * time.Second)
		closePublisher(con, ch)
		consumer.Close()
		defer wg.Done()
	}()
	// Act
	err = app.RecommendationHandler.ProcessMessages()
	wg.Wait()

	// Assert
	suite.Suite.Assert().Nil(err)

	var actualRecommendation = models.Recommendation{}
	err = sqlLiteDB.Where("book_id = ?", 3).First(&actualRecommendation).Error
	suite.Suite.Assert().Nil(err)
	suite.Suite.Assert().Equal(2, actualRecommendation.Points)

	var actualRecommendation2 = models.Recommendation{}
	err = sqlLiteDB.Where("book_id = ?", 4).First(&actualRecommendation2).Error
	suite.Suite.Assert().Nil(err)
	suite.Suite.Assert().Equal(7, actualRecommendation2.Points)

	// clean up
	sqlLiteDB.Unscoped().Where("book_id = ?", "3").Delete(&models.Recommendation{})
	sqlLiteDB.Unscoped().Where("book_id = ?", "4").Delete(&models.Recommendation{})
}
func (suite *IntegrationSuite) TestRabbitMQ_ShouldLogErrorMessageIfCannotConnectToRabbitMQ() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not config: %v", err)
	}
	loggerMock := &LoggerMock{}

	// Assert
	loggerMock.On("LogError", mock.MatchedBy(func(err error) bool {
		return strings.Contains(err.Error(), "failed to connect to RabbitMQ")
	})).Once()

	//act
	_ = application.InitializeRabbitMqConsumer(config, loggerMock)
}

func newPublisher(url, queueName string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	//defer ch.Close()

	// Declare the queue
	_, err = ch.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}
	return conn, ch
}

func publish(ch *amqp.Channel, queueName string, message interface{}) error {

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Publish the message
	err = ch.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	log.Printf("[x] Sent to queue %s: %s\n", queueName, body)
	return nil
}

func closePublisher(conn *amqp.Connection, ch *amqp.Channel) {
	if ch != nil {
		ch.Close()
	}
	if conn != nil {
		conn.Close()
	}

}
