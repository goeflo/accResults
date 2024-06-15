package data

type ResultData struct {
	SessionType   string        `json:"sessionType"`
	TrackName     string        `json:"trackName"`
	ServerName    string        `json:"serverName"`
	SessionIndex  uint          `json:"sessionIndex"`
	SessionResult SessionResult `json:"sessionResult"`
	Laps          []Lap         `json:"laps"`
	Penalties     []Penalty     `json:"penalties"`
}

type Penalty struct {
}

type Lap struct {
	CarID          uint   `json:"carId"`
	DriverIndex    uint   `json:"driverIndex"`
	Laptime        uint   `json:"laptime"`
	IsValidForBest bool   `json:"isValidForBest"`
	Splits         []uint `json:"splits"`
}

type SessionResult struct {
	BestLap          uint              `json:"bestlap"`
	BestSplits       []uint            `json:"bestSplits"`
	IsWetSession     uint              `json:"isWetSession"`
	Type             uint              `json:"type"`
	LeaderBoardLines []LeaderBoardLine `json:"leaderBoardLines"`
}

type LeaderBoardLine struct {
	Car                     Car       `json:"car"`
	CurrentDriver           Driver    `json:"currentDriver"`
	CurrentDriverIndex      uint      `json:"currentDriverIndex"`
	Timing                  Timing    `json:"timing"`
	MissingMandatoryPitstop uint      `json:"missingMandatoryPitstop"`
	DriverTotalTimes        []float32 `json:"driverTotalTimes"`
}

type Timing struct {
	LastLap     uint   `json:"lastLap"`
	LastSplits  []uint `json:"lastSplits"`
	BestLap     uint   `json:"bestLap"`
	BestSplits  []uint `json:"bestSplits"`
	TotalTime   uint   `json:"totalTime"`
	LapCount    uint   `json:"lapCount"`
	LastSplitID uint   `json:"lastSplitId"`
}

type Car struct {
	CarID       uint     `json:"carId"`
	RaceNumber  uint     `json:"raceNumber"`
	CarModel    uint     `json:"carModel"`
	CupCategory uint     `json:"cupCategory"`
	TeamName    string   `json:"teamName"`
	Nationality uint     `json:"nationality"`
	Drivers     []Driver `json:"drivers"`
}

type Driver struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	ShortName string `json:"shortName"`
	PlayerID  string `json:"playerId"`
}
