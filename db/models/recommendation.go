package models

import (
	"gorm.io/gorm"
)

type Recommendation struct {
	gorm.Model
	BookId int64 `gorm:"index"`
	Points int64
}
