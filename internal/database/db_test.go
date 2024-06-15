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
	if err := result.Read("../../test/191003_235558_R.json"); err != nil {
		t.Error(err)
	}

	db := NewSqlLite("test.db")
	fmt.Printf("db filename: %v\n", db.filename)

	if err := db.NewResult(result); err != nil {
		t.Error(err)
	}

	os.Remove("test.db")

}
