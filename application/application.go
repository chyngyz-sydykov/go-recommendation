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
	"gorm.io/gorm"
)

type App struct {
	RecommendationHandler *handlers.RecommendationHandler
	DB                    *gorm.DB
	MessageBrokerConsumer messagebroker.MessageBrokerConsumerInterface
}

func InitializeApplication() *App {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not config: %v", err)
	}

	db := initializeDatabase()

	consumer := InitializeRabbitMqConsumer(config)

	recommendationService := recommendation.NewRecommendationService(db)

	logger := logger.NewLogger()
	commonHandler := handlers.NewCommonHandler(logger)

	recommendationHandler := handlers.NewRecommendationHandler(commonHandler, consumer, recommendationService)
	app := &App{
		RecommendationHandler: recommendationHandler,
		DB:                    db,
		MessageBrokerConsumer: consumer,
	}
	return app
}
func (app *App) Start() {
	app.RecommendationHandler.ProcessMessages()
}
func (app *App) ShutDown() {

	fmt.Println("Application exited gracefully.")
	app.MessageBrokerConsumer.Close()

	db, err := app.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	if db != nil {
		db.Close()
	}
}

func InitializeRabbitMqConsumer(config *config.Config) messagebroker.MessageBrokerConsumerInterface {
	rabbitMQURL := "amqp://" + config.RabbitMqUser + ":" + config.RabbitMqPassword + "@" + config.RabbitMqContainerName + ":5672/"
	consumer, err := messagebroker.NewRabbitMQConsumer(rabbitMQURL, config.RabbitMqQueueName)
	if err != nil {
		log.Fatalf("Failed to initialize message publisher: %v", err)
	}

	return consumer
}

func initializeDatabase() *gorm.DB {
	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		log.Fatalf("Could not load database config: %v", err)
	}
	dbInstance, err := db.InitializeDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Coult not initialize db connection %v", err)
	}
	db.Migrate()
	return dbInstance

}
