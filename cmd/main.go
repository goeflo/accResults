package main

import (
	"accResults/config"
	"accResults/internal/database"
	"accResults/server"
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
)

func main() {

	// command line flag for configuration file
	configfile := flag.String("config", "", "path to the config toml file")
	flag.Parse()

	if *configfile == "" {
		log.Println("config file not specified (--config)")
		os.Exit(-1)
	}

	config, err := config.NewRaceResultsConfig(*configfile)
	if err != nil {
		panic(err)
	}

	// create data directory if not exists
	if _, err := os.Stat(config.DataDir); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(config.DataDir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	// create sqlite database
	db := database.NewSqlLite(filepath.Join(config.DataDir, config.SqlLiteFilename))
	//db.Import(data.Result{})

	// start html server
	server := server.NewServer(config, db)
	server.Run()

}
