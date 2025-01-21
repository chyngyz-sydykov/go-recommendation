package main

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/chyngyz-sydykov/go-recommendation/application/handlers"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/db/models"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/logger"
	"github.com/chyngyz-sydykov/go-recommendation/internal/recommendation"
	"github.com/stretchr/testify/mock"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (suite *IntegrationSuite) TestRecommendationUpsert_ShouldSuccessfullySave() {
	testCases := []struct {
		name            string
		BookId          int
		Event           string
		ResultingPoints int
	}{

		{"InitialBookUpdate", 1, "bookUpdated", 1},
		{"RecurringBookUpdateOnUpdatedBook", 1, "bookUpdated", 2},
		{"InitialBookRated", 2, "bookRated", 3},
		{"RecurringBookUpdateOnRatedBook", 2, "bookUpdated", 4},
		{"RecurringBookRateOnRatedBook", 2, "bookRated", 7},
	}
	sqlLiteDB := initializeSqlLite()

	logger := logger.NewLogger()
	commonHandler := handlers.NewCommonHandler(logger)

	for _, testCase := range testCases {
		suite.T().Run(testCase.name, func(t *testing.T) {
			// assign
			bookMessage := recommendation.BookMessage{
				BookId: testCase.BookId,
				Event:  testCase.Event,
			}
			messageBody, _ := json.Marshal(bookMessage)

			mockDelivery := make(chan amqp.Delivery, 1)
			mockDelivery <- amqp.Delivery{
				Body: messageBody,
			}
			close(mockDelivery)

			readOnlyChan := (<-chan amqp.Delivery)(mockDelivery)

			var messageBrokerMock MessageBrokerMock

			messageBrokerMock.On("Consume").Return(readOnlyChan, nil)

			app := provideDependencies(suite, &messageBrokerMock, commonHandler)

			// Act
			err := app.RecommendationHandler.ProcessMessages()

			// Assert
			suite.Suite.Assert().Nil(err)

			var actualRecommendation = models.Recommendation{}
			err = sqlLiteDB.Where("book_id = ?", testCase.BookId).First(&actualRecommendation).Error
			suite.Suite.Assert().Nil(err)
			suite.Suite.Assert().Equal(testCase.ResultingPoints, actualRecommendation.Points)

		})
	}

	// clean up
	sqlLiteDB.Unscoped().Where("book_id = ?", "1").Delete(&models.Recommendation{})
	sqlLiteDB.Unscoped().Where("book_id = ?", "2").Delete(&models.Recommendation{})
}

func (suite *IntegrationSuite) TestMessageProccessing_ShouldLogErrorIfMessageBorkerCannotConsumeMessages() {
	// Arrange
	var messageBrokerMock MessageBrokerMock
	commonHandlerMock := &CommonHandlerMock{}

	messageBrokerMock.On("Consume").Return((<-chan amqp.Delivery)(nil), errors.New("consume error"))

	// Assert
	commonHandlerMock.On("HandleError", mock.MatchedBy(func(err error) bool {
		return strings.Contains(err.Error(), "failed to start consuming messages")
	})).Once()

	app := provideDependencies(suite, &messageBrokerMock, commonHandlerMock)

	// Act
	_ = app.RecommendationHandler.ProcessMessages()
}
func (suite *IntegrationSuite) TestMessageProccessing_ShouldLogErrorIfMessageCannotBeMarshalledToDto() {
	// Arrange
	var messageBrokerMock MessageBrokerMock
	commonHandlerMock := &CommonHandlerMock{}

	mockDelivery := make(chan amqp.Delivery, 1)
	mockDelivery <- amqp.Delivery{
		Body: []byte("invalid json"),
	}
	close(mockDelivery)

	readOnlyChan := (<-chan amqp.Delivery)(mockDelivery)

	messageBrokerMock.On("Consume").Return(readOnlyChan, nil)

	// Assert
	commonHandlerMock.On("HandleError", mock.MatchedBy(func(err error) bool {
		return strings.Contains(err.Error(), "failed to unmarshal message")
	})).Once()

	app := provideDependencies(suite, &messageBrokerMock, commonHandlerMock)

	// Act
	_ = app.RecommendationHandler.ProcessMessages()
}
func (suite *IntegrationSuite) TestMessageProccessing_ShouldLogErrorIfNotSupportedEventProvided() {
	// Arrange
	var messageBrokerMock MessageBrokerMock
	commonHandlerMock := &CommonHandlerMock{}

	bookMessage := recommendation.BookMessage{
		BookId: 1,
		Event:  "notSupportedEvent",
	}
	messageBody, _ := json.Marshal(bookMessage)

	mockDelivery := make(chan amqp.Delivery, 1)
	mockDelivery <- amqp.Delivery{
		Body: messageBody,
	}
	close(mockDelivery)

	readOnlyChan := (<-chan amqp.Delivery)(mockDelivery)

	messageBrokerMock.On("Consume").Return(readOnlyChan, nil)

	// Assert
	commonHandlerMock.On("HandleError", mock.MatchedBy(func(err error) bool {
		return strings.Contains(err.Error(), "point cannot be generated for following event")
	})).Once()

	app := provideDependencies(suite, &messageBrokerMock, commonHandlerMock)

	// Act
	_ = app.RecommendationHandler.ProcessMessages()
}

func (suite *IntegrationSuite) TestMessageProccessing_ShouldLogErrorIfCannotSaveToDatabase() {
	// Arrange
	var messageBrokerMock MessageBrokerMock
	commonHandlerMock := &CommonHandlerMock{}

	mockDelivery := make(chan amqp.Delivery, 1)
	mockDelivery <- amqp.Delivery{
		Body: []byte("{\"unknown_key\": 1}"),
	}
	close(mockDelivery)

	readOnlyChan := (<-chan amqp.Delivery)(mockDelivery)

	messageBrokerMock.On("Consume").Return(readOnlyChan, nil)

	// Assert
	commonHandlerMock.On("HandleError", mock.MatchedBy(func(err error) bool {
		return strings.Contains(err.Error(), "failed to create recommendation")
	})).Once()

	app := provideDependencies(suite, &messageBrokerMock, commonHandlerMock)

	// Act
	_ = app.RecommendationHandler.ProcessMessages()
}
