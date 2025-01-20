package recommendation

import (
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/db"
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/db/models"
)

type RecommendationServiceInterface interface {
	ProcessMessage(recommendation *RecommendationDTO) error
}

type RecommendationService struct {
	repository RecommendationRepository
	calculator PointCalculator
}

func NewRecommendationService(db db.DatabaseInterface) *RecommendationService {
	repository := NewRecommendationRepository(db)
	calculator := NewPointCalculator()
	return &RecommendationService{
		repository: *repository,
		calculator: *calculator,
	}
}

func (service *RecommendationService) ProcessMessage(RecommendationDTO *RecommendationDTO) error {
	err := service.assignPointByEventName(RecommendationDTO)
	if err != nil {
		return err
	}
	recommendation := service.mapToGorm(RecommendationDTO)
	return service.repository.Upsert(recommendation)
}

func (service *RecommendationService) assignPointByEventName(RecommendationDTO *RecommendationDTO) error {
	points, err := service.calculator.GetPoint(RecommendationDTO.Event)
	if err != nil {
		return err
	}

	RecommendationDTO.Points = points
	return nil
}

func (service *RecommendationService) mapToGorm(RecommendationDTO *RecommendationDTO) *models.Recommendation {
	var recommendation models.Recommendation
	recommendation.BookId = RecommendationDTO.BookId
	recommendation.Points = RecommendationDTO.Points
	return &recommendation
}
