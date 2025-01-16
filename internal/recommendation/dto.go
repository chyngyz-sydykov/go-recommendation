package recommendation

import (
	"time"
)

type RecommendationDTO struct {
	ID     int `json:"id"`
	BookId int `json:"book_id"`
	Points int `json:"points"`
}

type BookMessage struct {
	ID       uint      `json:"id"`
	Title    string    `json:"title"`
	ICBN     string    `json:"icbn"`
	EditedAt time.Time `json:"EditedAt"`
	Event    string    `json:"event"`
}
