package models

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	Firstname string
	Lastname  string
}
