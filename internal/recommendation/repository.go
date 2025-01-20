package recommendation

import (
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/db"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/db/models"
)

type RecommendationRepository struct {
	db db.DatabaseInterface
}

func NewRecommendationRepository(db db.DatabaseInterface) *RecommendationRepository {
	return &RecommendationRepository{db: db}
}

func (repository *RecommendationRepository) Upsert(recommendation *models.Recommendation) error {
	err := repository.db.Upsert(recommendation)
	if err != nil {
		return err
	}
	return nil
}
