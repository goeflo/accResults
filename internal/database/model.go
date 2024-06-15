package database

import "gorm.io/gorm"

type ResultFile struct {
	gorm.Model
	Filename string
}

type Race struct {
	gorm.Model
	ResultFileID uint
	Track        string
	ServerName   string
	SessionType  string
}
