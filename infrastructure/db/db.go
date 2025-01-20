package db

import "github.com/chyngyz-sydykov/go-recommendation/infrastructure/db/models"

type DatabaseInterface interface {
	Upsert(recommendation *models.Recommendation) error
	// Create(model any) error
	// Update(book *models.Recommendation, payload models.Recommendation) error
	Close() error
	Migrate()
}
