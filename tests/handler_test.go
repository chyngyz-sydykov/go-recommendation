package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/chyngyz-sydykov/go-recommendation/application"
	"github.com/chyngyz-sydykov/go-recommendation/application/handlers"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/config"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/db"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/messagebroker"
	"github.com/chyngyz-sydykov/go-recommendation/internal/recommendation"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type IntegrationSuite struct {
	suite.Suite
	db db.DatabaseInterface
}

func (suite *IntegrationSuite) SetupSuite() {
	db := initializeDatabase()
	suite.db = db
}

func (suite *IntegrationSuite) TearDownSuite() {
	suite.db.Close()
}

// TestSuite runs the test suite.
func TestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}

func initializeDatabase() db.DatabaseInterface {
	dbConfig, err := config.LoadSqLiteDBConfig()
	if err != nil {
		log.Fatalf("Could not config: %v", err)
	}

	// Connect to the target database using Gorm
	dbInstance, err := db.InitializeSqlLite(dbConfig)
	if err != nil {
		log.Fatalf("Coult not initialize db connection %v", err)
	}
	dbInstance.Migrate()
	return dbInstance
}

type MessageBrokerMock struct {
	mock.Mock
}

func (m *MessageBrokerMock) Consume() (<-chan amqp.Delivery, error) {
	args := m.Called()
	return args.Get(0).(<-chan amqp.Delivery), args.Error(1)
}

func (m *MessageBrokerMock) InitializeMessageBroker() {
	m.Called()
}

func (m *MessageBrokerMock) Close() {
}

type CommonHandlerMock struct {
	mock.Mock
}

func (m *CommonHandlerMock) HandleError(err error) {
	m.Called(err)
}

type LoggerMock struct {
	mock.Mock
}

func (m *LoggerMock) LogError(err error) {
	m.Called(err)
}

func provideDependencies(suite *IntegrationSuite, consumerMock messagebroker.MessageBrokerConsumerInterface, commonHandler handlers.CommonHandlerInterface) *application.App {

	recommendationService := recommendation.NewRecommendationService(suite.db)
	recommendationHandler := handlers.NewRecommendationHandler(commonHandler, consumerMock, recommendationService)

	app := &application.App{
		RecommendationHandler: recommendationHandler,
	}
	return app
}

func initializeSqlLite() *gorm.DB {
	dbConfig, err := config.LoadSqLiteDBConfig()
	if err != nil {
		log.Fatalf("Could not config: %v", err)
	}

	fullPath := fmt.Sprintf("%s/%s", dbConfig.Path, dbConfig.Name)
	db, err := gorm.Open(sqlite.Open(fullPath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	return db
}
