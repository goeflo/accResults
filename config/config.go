package config

import (
	"log/slog"

	"github.com/BurntSushi/toml"
)

type RaceResultConfig struct {
	HttpServerAddr  string `toml:"http_server_addr"`
	SqlLiteFilename string `toml:"sqlite_filename"`
	DataDir         string `toml:"data_dir"`
}

func NewRaceResultsConfig(filename string) (*RaceResultConfig, error) {
	slog.Info("reading configuration", "filename", filename)

	config := RaceResultConfig{}

	_, err := toml.DecodeFile(filename, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
