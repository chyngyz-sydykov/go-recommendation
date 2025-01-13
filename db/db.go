package db

import (
	"fmt"
	"log"

	"github.com/chyngyz-sydykov/go-web/infrastructure/config"
	"github.com/chyngyz-sydykov/go-web/infrastructure/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBInstance struct {
	*gorm.DB
}

var db *gorm.DB

func InitializeDatabase(dbConfig *config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Name, dbConfig.Port)
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	return db, nil
}

func Migrate() {
	err := db.AutoMigrate(&models.Author{}, &models.Book{})
	if err != nil {
		log.Fatal("failed to run migration:", err)
	}
	log.Println("Migration completed successfully.")
}
