package models

import (
	"gorm.io/gorm"
)

type Recommendation struct {
	gorm.Model
	BookId int `gorm:"index"`
	Points int
}
