package model

type Leaderboard struct {
	DriverID  uint
	CarID     uint
	Pos       uint
	No        uint
	Driver    string
	Vehicle   string
	Laps      uint
	Bestlap   string
	Totaltime string
	Gap       string
}

type Lap struct {
	Lap     uint
	Time    string
	Sector1 string
	Sector2 string
	Sector3 string
}

type Driver struct {
	Firstname  string
	Lastname   string
	Shortname  string
	Vehicle    string
	FastestLap uint
	Laps       []Lap
}

/*
Pos 	No 	Driver 	Vehicle 	Laps 	Best lap 	Gap
*/
