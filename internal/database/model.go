package database

import "gorm.io/gorm"

type ResultFiles struct {
	gorm.Model
	Filename string
}
