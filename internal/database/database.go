package database

import (
	"github.com/goeflo/accResults/internal/data"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database interface {
	NewResult(result data.Result) error
	GetRaces() ([]Race, error)
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
	db.AutoMigrate(
		&ResultFile{},
		&Race{},
	)
	return SqlLite{filename: filename, db: db}

}

func (s SqlLite) NewResult(data data.Result) error {
	resultFile := ResultFile{Filename: data.Filename}

	if result := s.db.Create(&resultFile); result.Error != nil {
		return result.Error
	}

	race := Race{
		ResultFileID: resultFile.ID,
		Track:        data.ResultData.TrackName,
		ServerName:   data.ResultData.ServerName,
		SessionType:  data.ResultData.SessionType,
	}

	if result := s.db.Create(&race); result.Error != nil {
		return result.Error
	}

	return nil
}

func (s SqlLite) GetRaces() ([]Race, error) {

	races := []Race{}
	if result := s.db.Order("ID desc").Find(&races); result.Error != nil {
		return nil, result.Error
	}

	return races, nil
}
