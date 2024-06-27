package database

import (
	"os"
	"testing"

	"github.com/goeflo/accResults/internal/data"
	"github.com/stretchr/testify/assert"
)

var db SqlLite
var result data.Result

func TestMain(m *testing.M) {
	db = NewSqlLite("test.db")
	defer os.Remove("test.db")

	result = data.NewResult()
	result.Read("../../test/191003_235558_R.json")
	db.NewResult(result)

	os.Exit(m.Run())

}
func TestDb(t *testing.T) {

	db.db.Create(&ResultFile{Filename: "kermit"})
	db.db.Create(&ResultFile{Filename: "piggy"})

	resultFile := ResultFile{}
	result := db.db.First(&resultFile, 2)
	assert.Nil(t, result.Error)
	assert.Equal(t, "kermit", resultFile.Filename)

}

func TestLeaderboard(t *testing.T) {

	lb, err := db.GetLeaderboard(1)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(lb))
	assert.Equal(t, result.ResultData.SessionResult.LeaderBoardLines[0].Timing.TotalTime, lb[0].Totaltime)

}
func TestCreateRace(t *testing.T) {
	cars, err := db.GetCarsForRace(1)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(cars))
}
