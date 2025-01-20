package recommendation

import "fmt"

type PointCalculator struct {
	BookUpdated int
	BookRated   int
}

var pointsMap = map[string]int{
	"bookUpdated": 1,
	"BookRated":   3,
}

func NewPointCalculator() *PointCalculator {
	return &PointCalculator{}
}

func (pc *PointCalculator) GetPoint(eventName string) (int, error) {
	// Check if the eventName exists in the map
	if point, ok := pointsMap[eventName]; ok {
		return point, nil
	}
	// Return an error if the eventName is not found
	return 0, fmt.Errorf("point cannot be generated for following event '%s'", eventName)
}
