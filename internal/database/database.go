package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database interface {
	NewResult(filename string) error
}

type SqlLite struct {
	filename string
	db       *gorm.DB
}

func NewSqlLite(filename string) SqlLite {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&ResultFiles{})
	return SqlLite{filename: filename, db: db}

}

func (s SqlLite) NewResult(filename string) error {
	r := ResultFiles{Filename: filename}

	if result := s.db.Create(&r); result.Error != nil {
		return result.Error
	}

	return nil
}
