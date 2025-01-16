package recommendation

import (
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/db/models"

	"gorm.io/gorm"
)

type RecommendationRepositoryInterface interface {
	Create(recommendation models.Recommendation) error
}

type RecommendationRepository struct {
	db *gorm.DB
}

func NewRecommendationRepository(db *gorm.DB) *RecommendationRepository {
	return &RecommendationRepository{db: db}
}

func (repository *RecommendationRepository) Create(recommendation *models.Recommendation) error {
	return repository.db.Create(&recommendation).Error
}
