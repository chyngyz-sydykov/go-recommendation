package db

import (
	"fmt"
	"log"

	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/config"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/db/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SqlLite struct {
	db *gorm.DB
}

func InitializeSqlLite(dbConfig *config.SqLiteDBConfig) (DatabaseInterface, error) {
	fullPath := fmt.Sprintf("%s/%s", dbConfig.Path, dbConfig.Name)
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

func (sqlite *SqlLite) Upsert(recommendation *models.Recommendation) error {
	err := sqlite.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "book_id"}}, // Conflict on `book_id`
			DoUpdates: clause.Assignments(map[string]interface{}{
				"points": gorm.Expr("points + ?", recommendation.Points), // Increment points
			}),
		}).Create(&recommendation).Error

		if err != nil {
			return err
		}

		return nil // Commit transaction
	})

	if err != nil {
		return err
	}
	return nil
}
func (sqlite *SqlLite) Update(book *models.Recommendation, payload models.Recommendation) error {
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
