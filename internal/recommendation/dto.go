package recommendation

import (
	"time"
)

type RecommendationDTO struct {
	ID     int    `json:"id"`
	BookId int    `json:"book_id"`
	Points int    `json:"points"`
	Event  string `json:"event"`
}

type BookMessage struct {
	BookId   int       `json:"book_id"`
	EditedAt time.Time `json:"EditedAt"`
	Event    string    `json:"event"`
}
