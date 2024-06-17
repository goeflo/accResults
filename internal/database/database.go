package database

import (
	"log/slog"

	"github.com/goeflo/accResults/internal/data"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database interface {
	NewResult(data data.Result) (uint, error)
	GetRaces() ([]Race, error)
	GetCars(raceID uint) ([]Car, error)
	AddDriver(driver Driver) error
	GetDriver(ID uint) (*Driver, error)
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
		&Car{},
		&LeaderBoard{},
		&Driver{},
		&Lap{},
	)
	return SqlLite{filename: filename, db: db}

}

func (s SqlLite) NewResult(data data.Result) (ID uint, err error) {
	resultFile := ResultFile{Filename: data.Filename}

	slog.Info("process new results", "serverName", data.ResultData.ServerName)

	if result := s.db.Create(&resultFile); result.Error != nil {
		return 0, result.Error
	}

	race := Race{
		ResultFileID: resultFile.ID,
		Track:        data.ResultData.TrackName,
		ServerName:   data.ResultData.ServerName,
		SessionType:  data.ResultData.SessionType,
	}

	if result := s.db.Create(&race); result.Error != nil {
		return 0, result.Error
	}

	slog.Info("leaderboard", "size", len(data.ResultData.SessionResult.LeaderBoardLines))

	for _, leaderboardline := range data.ResultData.SessionResult.LeaderBoardLines {
		car := Car{
			RaceID:         race.ID,
			AccResultCarID: leaderboardline.Car.CarID,
			CarModel:       leaderboardline.Car.CarModel,
			RaceNumber:     leaderboardline.Car.RaceNumber,
			TeamName:       leaderboardline.Car.TeamName,
			Nationality:    leaderboardline.Car.Nationality,
		}

		if result := s.db.Create(&car); result.Error != nil {
			slog.Error("error adding new car", "err", result.Error)
			return 0, result.Error
		}

		for _, driver := range leaderboardline.Car.Drivers {
			if err := s.AddDriver(Driver{
				CarID:     car.ID,
				FirstName: driver.FirstName,
				LastName:  driver.LastName,
				ShortName: driver.ShortName,
			}); err != nil {
				return 0, err
			}
		}

		lb := LeaderBoard{
			RaceID:      race.ID,
			CarID:       race.ID,
			LapCount:    leaderboardline.Timing.LapCount,
			LastLaptime: leaderboardline.Timing.LastLap,
			BestLaptime: leaderboardline.Timing.BestLap,
			Totaltime:   leaderboardline.Timing.TotalTime,
		}
		if leaderboardline.MissingMandatoryPitstop == 0 {
			lb.MissingPitstop = false
		} else {
			lb.MissingPitstop = true
		}

		if result := s.db.Create(&lb); result.Error != nil {
			slog.Error("error adding new leaderboard", "err", result.Error)
			return 0, result.Error
		}

	}
	return resultFile.ID, nil
}

func (s SqlLite) AddDriver(driver Driver) error {
	slog.Debug("add new driver", "shortName", driver.ShortName)
	if result := s.db.Create(&driver); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s SqlLite) GetDriver(ID uint) (driver *Driver, err error) {
	if result := s.db.First(driver, ID); result.Error != nil {
		return nil, result.Error
	}
	return driver, nil

}

func (s SqlLite) GetCars(raceID uint) (cars []Car, err error) {
	if result := s.db.Where(&Car{RaceID: raceID}).Find(&cars); result.Error != nil {
		return nil, result.Error
	}
	return cars, nil
}

func (s SqlLite) GetRaces() (races []Race, err error) {
	if result := s.db.Order("ID desc").Find(&races); result.Error != nil {
		return nil, result.Error
	}
	return races, nil
}
