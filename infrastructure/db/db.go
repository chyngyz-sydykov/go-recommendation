package db

import "github.com/chyngyz-sydykov/go-recommendation/infrastructure/db/models"

type DatabaseInterface interface {
	Upsert(recommendation *models.Recommendation) error
	Close() error
	Migrate()
}
