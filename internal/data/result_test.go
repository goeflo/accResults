package data

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {

	r := NewResult()
	if err := r.Read("../../test/191003_235558_R.json"); err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(r.ResultData.Laps), 65)
	assert.Equal(t, r.ResultData.SessionResult.LeaderBoardLines[0].CurrentDriver.PlayerID, "123")

	fmt.Printf("laptime: %v\n", r.ResultData.Laps[0].Laptime)
	//fmt.Printf("laptime: %v\n", r.ResultData.Laps[0].Laptime/1000)

	fmt.Printf("%v\n", convertMilliseconds(r.ResultData.Laps[0].Laptime))

}
