package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDb(t *testing.T) {

	db := NewSqlLite("test")
	fmt.Printf("db filename: %v\n", db.filename)

	db.db.Create(&ResultFiles{Filename: "kermit"})
	db.db.Create(&ResultFiles{Filename: "piggy"})

	resultFile := ResultFiles{}
	result := db.db.First(&resultFile, 1)
	assert.Nil(t, result.Error)
	assert.Equal(t, "kermit", resultFile.Filename)

	os.Remove("test")
}
