package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Result struct {
	Filename   string
	ResultData ResultData
}

func NewResult() Result {
	return Result{}
}

func (r *Result) Store(dataDir string, data []byte) error {

	weekDirName := r.createDataWeekDirname(dataDir)
	if _, err := os.Stat(weekDirName); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(weekDirName, os.ModePerm); err != nil {
			panic(err)
		}
	}

	r.Filename = filepath.Join(weekDirName, uuid.New().String()+".json")

	slog.Info("save result", "filename", r.Filename)
	if err := os.WriteFile(r.Filename, data, 0644); err != nil {
		return err
	}

	// if err := json.Unmarshal(data, &r.ResultData); err != nil {
	// 	return err
	// }

	return nil
}

func (r *Result) createDataWeekDirname(dataDir string) string {
	tn := time.Now().UTC()
	year, week := tn.ISOWeek()
	return filepath.Join(dataDir, fmt.Sprintf("%v_%v", year, week))
}

// Read read result data from file set in Result.Filename
// json result is stored in UTF16-LE format by acc.
func (r *Result) Read(filename string) error {

	r.Filename = filename

	file, err := os.Open(r.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	winutf := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	decoder := winutf.NewDecoder()
	reader := transform.NewReader(file, unicode.BOMOverride(decoder))

	raw, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw, &r.ResultData); err != nil {
		return err
	}

	return nil
}

func ConvertMilliseconds(milliseconds uint) string {

	d := time.Duration(milliseconds) * time.Millisecond

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	millis := int(d.Milliseconds()) % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, millis)

}
