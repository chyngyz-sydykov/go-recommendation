package recommendation

import (
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/db"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/db/models"
)

type RecommendationServiceInterface interface {
	Create(recommendation *RecommendationDTO) error
}

type RecommendationService struct {
	repository RecommendationRepository
}

func NewRecommendationService(db db.DatabaseInterface) *RecommendationService {
	repository := NewRecommendationRepository(db)
	return &RecommendationService{
		repository: *repository,
	}
}

func (service *RecommendationService) Create(RecommendationDTO *RecommendationDTO) error {
	recommendation := service.mapToGorm(RecommendationDTO)
	return service.repository.Create(recommendation)
}

func (service *RecommendationService) mapToGorm(RecommendationDTO *RecommendationDTO) *models.Recommendation {
	var recommendation models.Recommendation
	recommendation.BookId = RecommendationDTO.BookId
	recommendation.Points = RecommendationDTO.Points
	return &recommendation
}
