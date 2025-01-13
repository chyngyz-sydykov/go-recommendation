package application

import (
	"log"
	"time"

	"github.com/chyngyz-sydykov/go-recommendation/application/handlers"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/config"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/db"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/logger"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/messagebroker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

type App struct {
}

func InitializeApplication() *App {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not config: %v", err)
	}

	db := initializeDatabase()
	//defer db.Close()

	rabbitMqPublisher := initializeRabbitMqPublisher(config)
	//defer rabbitMqPublisher.Close()

	logger := logger.NewLogger()
	commonHandler := handlers.NewCommonHandler(logger)

	ratingClient := initializeRatingGrpcClient()
	ratingService := rating.NewRatingService(ratingClient, time.Duration(config.GrpcTimeoutDuration)*time.Second)
	ratingHandler := handlers.NewRatingHandler(ratingService, *commonHandler)

	bookService := book.NewBookService(db, rabbitMqPublisher, ratingService)
	bookHandler := handlers.NewBookHandler(*bookService, *commonHandler)

	app := &App{
		BookHandler:   *bookHandler,
		RatingHandler: *ratingHandler,
	}
	return app
}

func initializeRatingGrpcClient() pb.RatingServiceClient {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not config: %v", err)
	}

	conn, err := grpc.NewClient(config.RatingServiceServer+":"+config.RatingServicePort, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	//defer conn.Close()

	// Create the gRPC client for the RatingService.
	ratingClient := pb.NewRatingServiceClient(conn)
	return ratingClient
}
func initializeRabbitMqPublisher(config *config.Config) messagebroker.MessageBrokerInterface {
	rabbitMQURL := "amqp://" + config.RabbitMqUser + ":" + config.RabbitMqPassword + "@" + config.RabbitMqContainerName + ":5672/"
	publisher, err := messagebroker.NewRabbitMQPublisher(rabbitMQURL, config.RabbitMqQueueName)

	if err != nil {
		log.Fatalf("Failed to initialize message publisher: %v", err)
	}
	publisher.InitializeMessageBroker()
	return publisher
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
