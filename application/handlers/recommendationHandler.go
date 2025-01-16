package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/messagebroker"
	"github.com/chyngyz-sydykov/go-recommendation/internal/recommendation"
)

type RecommendationHandler struct {
	service               recommendation.RecommendationServiceInterface
	MessageBrokerConsumer messagebroker.MessageBrokerConsumerInterface
	commonHandler         CommonHandler
}

func NewRecommendationHandler(commonHandler CommonHandler, consumer messagebroker.MessageBrokerConsumerInterface, service recommendation.RecommendationServiceInterface) *RecommendationHandler {
	return &RecommendationHandler{
		service:               service,
		MessageBrokerConsumer: consumer,
		commonHandler:         commonHandler,
	}
}
func (handler *RecommendationHandler) ProcessMessages() error {
	msgs, err := handler.MessageBrokerConsumer.Consume()

	if err != nil {
		tempError := fmt.Errorf("failed to start consuming messages: %w", err)
		handler.commonHandler.HandleError(tempError)
	}
	for msg := range msgs {
		fmt.Printf("Message: %s\n", msg.Body)
		var bookMessage recommendation.BookMessage
		err := json.Unmarshal(msg.Body, &bookMessage)
		if err != nil {
			tempError := fmt.Errorf("failed to unmarshal message: %w", err)
			handler.commonHandler.HandleError(tempError)
			msg.Reject(false)
			continue
		}
		recommendation := &recommendation.RecommendationDTO{
			BookId: int(bookMessage.ID),
			Points: 1,
		}

		err = handler.service.Create(recommendation)
		if err != nil {
			tempError := fmt.Errorf("failed to create recommendation: %w", err)
			handler.commonHandler.HandleError(tempError)
			return tempError
		}

		msg.Ack(false)
	}
	return nil
}
