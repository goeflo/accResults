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

type Car struct {
	gorm.Model
	RaceID         uint `gorm:"index"`
	RaceNumber     uint
	CarModel       uint
	CupCategory    uint
	TeamName       string
	Nationality    uint
	AccResultCarID uint `gorm:"index"`
}

type LeaderBoard struct {
	gorm.Model
	RaceID uint `gorm:"index"`
	CarID  uint
	//DriverID       uint
	LapCount       uint
	LastLaptime    uint
	BestLaptime    uint
	Totaltime      uint
	MissingPitstop bool
}

type Driver struct {
	gorm.Model
	CarID          uint `gorm:"index"`
	FirstName      string
	LastName       string
	ShortName      string
	PlayerID       string
	LapTimeAverage uint
}

type Lap struct {
	gorm.Model
	CarID            uint
	DriverID         uint
	Laptime          uint
	IsValid          bool
	FastestLapInRace bool
	Splits           string
	Split1           uint
	Split2           uint
	Split3           uint
}
