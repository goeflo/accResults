package views

type Leaderboard struct {
	DriverID uint
	CarID    uint
	Pos      uint
	No       uint
	Driver   string
	Vehicle  string
	Laps     uint
	Bestlap  string
	Gap      string
}

/*
Pos 	No 	Driver 	Vehicle 	Laps 	Best lap 	Gap
*/
