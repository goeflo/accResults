package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/goeflo/accResults/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestDb(t *testing.T) {

	db := NewSqlLite("test.db")
	fmt.Printf("db filename: %v\n", db.filename)

	db.db.Create(&ResultFile{Filename: "kermit"})
	db.db.Create(&ResultFile{Filename: "piggy"})

	resultFile := ResultFile{}
	result := db.db.First(&resultFile, 1)
	assert.Nil(t, result.Error)
	assert.Equal(t, "kermit", resultFile.Filename)

	os.Remove("test.db")
}

func TestCreateRace(t *testing.T) {
	result := data.NewResult()
	err := result.Read("../../test/191003_235558_R.json")
	assert.Nil(t, err)

	db := NewSqlLite("test.db")

	resultID, err := db.NewResult(result)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, resultID)

	cars, err := db.GetCars(1)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(cars))

	//	fmt.Printf("%+v\n", cars)

	os.Remove("test.db")

}
