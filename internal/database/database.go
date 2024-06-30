package database

import (
	"log"
	"math"

	"github.com/goeflo/accResults/internal/data"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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
	log.Printf("adding new results ...")

	accCarIDToDriverID := make(map[uint]uint)
	accCarIDToDBCarID := make(map[uint]uint)
	averageLapTimeFroDriverID := make(map[uint][]uint)

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
			log.Printf("error adding new car %v, %v\n", car, result.Error)
			return 0, result.Error
		}

		accCarIDToDBCarID[leaderboardline.Car.CarID] = car.ID

		for _, driver := range leaderboardline.Car.Drivers {
			dr := &Driver{
				CarID:     car.ID,
				FirstName: driver.FirstName,
				LastName:  driver.LastName,
				ShortName: driver.ShortName,
			}

			if err := s.AddDriver(dr); err != nil {
				return 0, err
			}

			accCarIDToDriverID[car.AccResultCarID] = dr.ID
		}

		lb := LeaderBoard{
			RaceID:      race.ID,
			CarID:       car.ID,
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
			log.Printf("error adding new leaderboard %v, %v\n", lb, result.Error)
			return 0, result.Error
		}
	}

	fastestLap := uint(math.MaxInt)
	fastestLapIdx := uint(0)
	for _, lap := range data.ResultData.Laps {
		l := Lap{
			CarID:            accCarIDToDBCarID[lap.CarID],
			DriverID:         accCarIDToDriverID[lap.CarID],
			Laptime:          lap.Laptime,
			IsValid:          lap.IsValidForBest,
			FastestLapInRace: false,
			Split1:           lap.Splits[0],
			Split2:           lap.Splits[1],
			Split3:           lap.Splits[2],
		}
		if result := s.db.Create(&l); result.Error != nil {
			log.Printf("error adding new lap %v, %v\n", l, result.Error)
			return 0, result.Error
		}
		if lap.IsValidForBest && lap.Laptime < fastestLap {
			fastestLap = lap.Laptime
			fastestLapIdx = l.ID
		}

		if lap.IsValidForBest {
			avgLapTime := averageLapTimeFroDriverID[accCarIDToDriverID[lap.CarID]]
			avgLapTime = append(avgLapTime, lap.Laptime)
			averageLapTimeFroDriverID[accCarIDToDriverID[lap.CarID]] = avgLapTime
		}

	}

	for driverID, avgLapTime := range averageLapTimeFroDriverID {
		driver := Driver{}
		s.db.First(&driver, driverID)

		calculatedAcgLapTime := uint(0)
		for i := range avgLapTime {
			calculatedAcgLapTime += avgLapTime[i]
		}
		calculatedAcgLapTime = calculatedAcgLapTime / uint(len(avgLapTime))

		driver.LapTimeAverage = calculatedAcgLapTime
		if result := s.db.Save(&driver); result.Error != nil {
			log.Printf("error updating avg lapt time for driver %v %v\n", driver.ID, result.Error)
			return 0, result.Error
		}
	}

	if fastestLapIdx != 0 {
		lap := Lap{}
		s.db.First(&lap, fastestLapIdx)
		lap.FastestLapInRace = true

		if result := s.db.Save(&lap); result.Error != nil {
			log.Printf("error updating flastest lap lapIdx %v %v\n", fastestLapIdx, result.Error)
			return 0, result.Error
		}
	}

	return resultFile.ID, nil
}

func (s SqlLite) GetLeaderboard(raceID uint) (leaderbord []LeaderBoard, err error) {
	if result := s.db.Order("lap_count desc, totaltime").Where(&LeaderBoard{RaceID: raceID}).Find(&leaderbord); result.Error != nil {
		return nil, result.Error
	}
	for i := range leaderbord {
		log.Printf("%v laps:%v, total time:%v\n", i, leaderbord[i].LapCount, leaderbord[i].Totaltime)
	}
	return leaderbord, nil
}

func (s SqlLite) AddDriver(driver *Driver) error {
	if result := s.db.Create(driver); result.Error != nil {
		return result.Error
	}
	log.Printf("new driver ID:%v, Shortname:%v\n", driver.ID, driver.ShortName)
	return nil
}

func (s SqlLite) GetLapsForDriver(driverID uint) (laps []Lap, err error) {
	if result := s.db.Where(&Lap{DriverID: driverID}).Find(&laps); result.Error != nil {
		return nil, result.Error
	}
	return laps, nil
}

func (s SqlLite) GetDriver(ID uint) (driver *Driver, err error) {
	if result := s.db.First(driver, ID); result.Error != nil {
		return nil, result.Error
	}
	return driver, nil
}

func (s SqlLite) GetDriverOnCar(carID uint) (driver *Driver, err error) {
	driver = &Driver{}
	if result := s.db.Where(&Driver{CarID: carID}).Find(&driver); result.Error != nil {
		return nil, result.Error
	}
	return driver, nil

}

func (s SqlLite) GetCar(ID uint) (car *Car, err error) {
	car = &Car{}
	if result := s.db.First(car, ID); result.Error != nil {
		return nil, result.Error
	}
	return car, nil
}

func (s SqlLite) GetCarsForRace(raceID uint) (cars []Car, err error) {
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

func (s SqlLite) GetRace(ID uint) (race *Race, err error) {
	race = &Race{}
	if result := s.db.First(race, ID); result.Error != nil {
		return nil, result.Error
	}
	return race, nil
}
