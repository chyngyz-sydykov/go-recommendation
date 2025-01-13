package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ApplicationAddress     string
	ApplicationPort        string
	ApplicationEnvironment string
	RatingServicePort      string
	RatingServiceServer    string
	GrpcTimeoutDuration    int
	RabbitMqUser           string
	RabbitMqPassword       string
	RabbitMqQueueName      string
	RabbitMqContainerName  string
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
}

func LoadConfig() (*Config, error) {
	err := loadEnvFile()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	config := &Config{
		ApplicationAddress:     getEnv("APPLICATION_ADDRESS", "/"),
		ApplicationPort:        getEnv("APPLICATION_PORT", "1111"),
		ApplicationEnvironment: getEnv("APPLICATION_ENVIRONMENT", "local"),
		RatingServiceServer:    getEnv("RATING_SERVICE_SERVER", "localhost"),
		RatingServicePort:      getEnv("RATING_SERVICE_PORT", "1112"),
		GrpcTimeoutDuration:    getIntEnv("GRPC_TIMEOUT_DURATION", 30),
		RabbitMqUser:           getEnv("RABBITMQ_USER", "guest"),
		RabbitMqPassword:       getEnv("RABBITMQ_PASSWORD", "guest"),
		RabbitMqQueueName:      getEnv("RABBITMQ_QUEUE_NAME", "queue-name"),
		RabbitMqContainerName:  getEnv("RABBITMQ_CONTAINER_NAME", "go_web_rabbitmq"),
	}

	return config, nil
}

func LoadDBConfig() (*DBConfig, error) {
	err := loadEnvFile()
	if err != nil {
		return nil, fmt.Errorf("error loading env file")
	}

	dbConfig := &DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		Name:     getEnv("DB_DATABASE", "database_name"),
		Username: getEnv("DB_USERNAME", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
	}

	return dbConfig, nil
}
func loadEnvFile() error {
	rootDir := os.Getenv("ROOT_DIR")
	envFileName := rootDir + "/.env"

	if os.Getenv("APP_ENV") != "development" {
		envFileName = rootDir + "/.env." + os.Getenv("APP_ENV")
	}

	err := godotenv.Load(envFileName)
	return err
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		number, err := strconv.Atoi(value)
		if err != nil {
			fmt.Errorf("error cannot convert string %s to int. Returning default value", value)
			return defaultValue
		}
		return number
	}
	return defaultValue
}
