package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ApplicationEnvironment string
	RabbitMqUser           string
	RabbitMqPassword       string
	RabbitMqQueueName      string
	RabbitMqContainerName  string
}

type PostgreDBConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
}
type SqLiteDBConfig struct {
	Path string
	Name string
}

func LoadConfig() (*Config, error) {
	err := loadEnvFile()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	config := &Config{
		ApplicationEnvironment: getEnv("APPLICATION_ENVIRONMENT", "local"),
		RabbitMqUser:           getEnv("RABBITMQ_USER", "guest"),
		RabbitMqPassword:       getEnv("RABBITMQ_PASSWORD", "guest"),
		RabbitMqQueueName:      getEnv("RABBITMQ_QUEUE_NAME", "queue-name"),
		RabbitMqContainerName:  getEnv("RABBITMQ_CONTAINER_NAME", "go_web_rabbitmq"),
	}

	return config, nil
}

func LoadPostgreDBConfig() (*PostgreDBConfig, error) {
	err := loadEnvFile()
	if err != nil {
		return nil, fmt.Errorf("error loading env file")
	}

	dbConfig := &PostgreDBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		Name:     getEnv("DB_DATABASE", "database_name"),
		Username: getEnv("DB_USERNAME", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
	}

	return dbConfig, nil
}

func LoadSqLiteDBConfig() (*SqLiteDBConfig, error) {
	err := loadEnvFile()
	if err != nil {
		return nil, fmt.Errorf("error loading env file")
	}

	dbConfig := &SqLiteDBConfig{
		Path: getEnv("SQLITE_PATH", "/app/data"),
		Name: getEnv("SQLITE_DB_NAME", "tmp.sqlite3"),
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

// func getIntEnv(key string, defaultValue int) (int, error) {
// 	if value, exists := os.LookupEnv(key); exists {
// 		number, err := strconv.Atoi(value)
// 		if err != nil {
// 			return defaultValue, fmt.Errorf("error cannot convert string %s to int. Returning default value", value)
// 		}
// 		return number, nil
// 	}
// 	return defaultValue, nil
// }
