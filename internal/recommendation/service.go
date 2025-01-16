package recommendation

import (
	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/db/models"

	"gorm.io/gorm"
)

type RecommendationServiceInterface interface {
	Create(recommendation *RecommendationDTO) error
}

type RecommendationService struct {
	repository RecommendationRepository
}

func NewRecommendationService(db *gorm.DB) *RecommendationService {
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

// func (service *RecommendationService) publishMessage(RecommendationDTO *RecommendationDTO, event string) error {

// 	recommendationMessage := recommendationMessage{
// 		ID:       recommendation.ID,
// 		Title:    recommendation.Title,
// 		ICBN:     recommendation.ICBN,
// 		EditedAt: time.Now(),
// 		Event:    event,
// 	}

// 	if err := service.messageBroker.Publish(recommendationMessage); err != nil {
// 		log.Fatalf("Failed to publish event: %v", err)
// 	}
// 	return nil
// }
