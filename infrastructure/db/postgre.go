package db

import (
	"fmt"
	"log"

	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/config"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
}

func InitializePostgres(dbConfig *config.PostgreDBConfig) (DatabaseInterface, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Name, dbConfig.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	return &Postgres{
		db: db,
	}, nil

}

func (sqlite *Postgres) Migrate() {
	err := sqlite.db.AutoMigrate(&models.Recommendation{})
	if err != nil {
		log.Fatal("failed to run migration:", err)
	}
	log.Println("Migration completed successfully.")
}

func (sqlite *Postgres) Create(recommendation any) error {
	err := sqlite.db.AutoMigrate(&models.Recommendation{})
	if err != nil {
		log.Fatal("failed to run migration:", err)
		return err
	}
	log.Println("Migration completed successfully.")
	return nil
}

func (sqlite *Postgres) Close() error {
	db, err := sqlite.db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	if db != nil {
		db.Close()
	}
	return nil
}
