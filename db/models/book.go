package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title       string
	ICBN        string
	PublishedAt *time.Time
	AuthorId    int64 `gorm:"index"`
	Author      Author
}
