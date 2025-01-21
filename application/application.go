package application

import (
	"fmt"
	"log"

	"github.com/chyngyz-sydykov/go-recommendation/application/handlers"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/config"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/db"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/logger"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/messagebroker"
	"github.com/chyngyz-sydykov/go-recommendation/internal/recommendation"
)

type App struct {
	RecommendationHandler *handlers.RecommendationHandler
	DB                    db.DatabaseInterface
	MessageBrokerConsumer messagebroker.MessageBrokerConsumerInterface
}

func InitializeApplication() *App {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not config: %v", err)
	}

	logger := logger.NewLogger()

	sqlLite := initializeSqlLiteDatabase()

	consumer := InitializeRabbitMqConsumer(config, logger)

	recommendationService := recommendation.NewRecommendationService(sqlLite)

	commonHandler := handlers.NewCommonHandler(logger)

	recommendationHandler := handlers.NewRecommendationHandler(commonHandler, consumer, recommendationService)
	app := &App{
		RecommendationHandler: recommendationHandler,
		DB:                    sqlLite,
		MessageBrokerConsumer: consumer,
	}
	return app
}
func (app *App) Start() {
	fmt.Println("Start")
	app.RecommendationHandler.ProcessMessages()
}
func (app *App) ShutDown() {

	fmt.Println("Application exited gracefully.")
	app.MessageBrokerConsumer.Close()

	if app.DB != nil {
		app.DB.Close()
	}
}

func InitializeRabbitMqConsumer(config *config.Config, logger logger.LoggerInterface) messagebroker.MessageBrokerConsumerInterface {
	rabbitMQURL := "amqp://" + config.RabbitMqUser + ":" + config.RabbitMqPassword + "@" + config.RabbitMqContainerName + ":5672/"
	consumer, err := messagebroker.NewRabbitMQConsumer(rabbitMQURL, config.RabbitMqQueueName)
	if err != nil {
		err = fmt.Errorf("failed to initialize message publisher: %v", err)
		logger.LogError(err)
	}

	return consumer
}

func initializeSqlLiteDatabase() db.DatabaseInterface {
	fmt.Println("initializeSqlLiteDatabase")
	dbConfig, err := config.LoadSqLiteDBConfig()
	if err != nil {
		log.Fatalf("Could not load database config: %v", err)
	}
	sqlLite, err := db.InitializeSqlLite(dbConfig)
	if err != nil {
		log.Fatalf("Coult not initialize db connection %v", err)
	}
	sqlLite.Migrate()
	return sqlLite

}
