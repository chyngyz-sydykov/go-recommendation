package db

import (
	"fmt"
	"log"

	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/config"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/db/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqlLite struct {
	db *gorm.DB
}

func InitializeSqlLite(dbConfig *config.SqLiteDBConfig) (DatabaseInterface, error) {
	fullPath := fmt.Sprintf("%s/%s", dbConfig.Path, dbConfig.Name)

	fmt.Println("fullPath", fullPath)

	db, err := gorm.Open(sqlite.Open(fullPath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	return &SqlLite{db: db}, nil

}

func (sqlite *SqlLite) Migrate() {
	err := sqlite.db.AutoMigrate(&models.Recommendation{})
	if err != nil {
		log.Fatal("failed to run migration:", err)
	}
	log.Println("Migration completed successfully.")
}

func (sqlite *SqlLite) Create(recommendation any) error {
	err := sqlite.db.Create(recommendation).Error
	if err != nil {
		log.Fatal("failed to run migration:", err)
		return err
	}
	log.Println("Created successfully.")
	return nil
}

func (sqlite *SqlLite) Close() error {
	db, err := sqlite.db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	if db != nil {
		db.Close()
	}
	return nil
}
