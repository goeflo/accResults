package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/goeflo/accResults/config"
	"github.com/goeflo/accResults/internal/data"
	"github.com/goeflo/accResults/internal/database"
	"github.com/goeflo/accResults/server/model"
	"github.com/goeflo/accResults/server/views"
	"github.com/gorilla/mux"
)

type Server struct {
	DB     database.SqlLite
	Config *config.RaceResultConfig
	Router *mux.Router
}

const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB

func NewServer(config *config.RaceResultConfig, db database.SqlLite) Server {

	router := mux.NewRouter()

	router.HandleFunc("/", handleHome(db)).Methods("GET")
	router.HandleFunc("/upload", handleUploadForm).Methods("GET")
	router.HandleFunc("/upload", handleUpload(config, db)).Methods("POST")
	router.HandleFunc("/details/{raceID}", handleDetails(db)).Methods("POST")
	router.HandleFunc("/details/{raceID}/laps", handleDetailsLaps(db)).Methods("POST")

	// handle all static content
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	return Server{
		DB:     db,
		Config: config,
		Router: router,
	}
}

func (s Server) Run() {
	log.Printf("starting http server address %v", s.Config.HttpServerAddr)
	http.ListenAndServe(s.Config.HttpServerAddr, s.Router)
}

func handleDetailsLaps(db database.SqlLite) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handleDetailsLaps...")
		/*
			params := mux.Vars(r)
			raceID, err := strconv.Atoi(params["raceID"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Printf("%v\n", err)
				return
			}

			race, err := db.GetRace(uint(raceID))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Printf("error getting race %v\n", err)
				return
			}

			cars, err := db.GetCarsForRace(race.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Printf("error getting cars for race %v %v\n", race.ID, err)
				return
			}

			divers := []data.Driver{}
			for _, car := range cars {
				driver, err := db.GetDriverOnCar(car.ID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					log.Printf("error getting driver on car %v %v\n", car.ID, err)
					return
				}

				laps, err := db.GetLaps()

			}*/
		//views.MakeDetailsLaps(race).Render(r.Context(), w)

	}
}

func handleDetails(db database.SqlLite) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handleDetails ...")
		params := mux.Vars(r)
		raceID, err := strconv.Atoi(params["raceID"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("%v\n", err)
			return
		}

		race, err := db.GetRace(uint(raceID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("error getting race %v\n", err)
			return
		}

		lb, err := db.GetLeaderboard(uint(raceID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("error getting leaderboard %v\n", err)
			return
		}

		pageLbs := []model.Leaderboard{}
		for i, line := range lb {

			car, err := db.GetCar(line.CarID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Printf("error getting carID %v %v\n", line.CarID, err)
				return
			}

			driver, err := db.GetDriverOnCar(car.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Printf("error getting driverID %v %v\n", car.ID, err)
				return
			}

			gap := "-"
			if i > 0 {
				if lb[i].LapCount == lb[0].LapCount {
					gap = data.ConvertMilliseconds(lb[i].Totaltime - lb[i-1].Totaltime)
				} else if lb[i].LapCount < lb[0].LapCount {
					lapDiff := lb[0].LapCount - lb[i].LapCount
					if lapDiff == 1 {
						gap = fmt.Sprintf("+%d lap", lapDiff)
					} else {
						gap = fmt.Sprintf("+%d laps", lapDiff)
					}

				}
			}
			log.Printf("totaltime %v, gap %v\n", lb[i].Totaltime, gap)

			pageLbs = append(pageLbs, model.Leaderboard{
				CarID:     line.CarID,
				DriverID:  driver.ID,
				No:        car.RaceNumber,
				Pos:       uint(i + 1),
				Driver:    fmt.Sprintf("%s %s", driver.FirstName, driver.LastName),
				Laps:      line.LapCount,
				Totaltime: data.ConvertMilliseconds(line.Totaltime),
				Gap:       gap,
				Bestlap:   data.ConvertMilliseconds(line.BestLaptime),
				Vehicle:   data.Cars[car.CarModel].Name,
			})

		}
		log.Printf("race track %v\n", race.Track)
		views.MakeDetailsResultPage(race, pageLbs).Render(r.Context(), w)
	}
}

func handleUpload(c *config.RaceResultConfig, db database.SqlLite) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handleUpload ...")

		r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
		if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
			log.Printf("%v\n", err.Error())
			http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
			return
		}

		resultFile, _, err := r.FormFile("race_result")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resultFile.Close()

		b, err := io.ReadAll(resultFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result := data.NewResult()
		if err := result.Store(c.DataDir, b); err != nil {
			log.Printf("%v\n", err.Error())
			http.Error(w, "result file was not saved "+err.Error(), http.StatusInternalServerError)
		}

		if err := result.Read(result.Filename); err != nil {
			log.Printf("%v\n", err.Error())
			http.Error(w, "can not read result file "+err.Error(), http.StatusInternalServerError)
		}

		_, err = db.NewResult(result)
		if err != nil {
			log.Printf("%v\n", err.Error())
			http.Error(w, "can not import race result into db", http.StatusInternalServerError)
		}

		races, err := db.GetRaces()
		if err != nil {
			log.Printf("%v\n", err.Error())
			http.Error(w, "can not read races data from db", http.StatusInternalServerError)
		}
		// TODO show upload result
		views.MakeIndex(races).Render(r.Context(), w)

	}
}

func handleHome(db database.SqlLite) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handleHome ...")

		races, err := db.GetRaces()
		if err != nil {
			log.Printf("%v\n", err.Error())
			http.Error(w, "can not read races data from db", http.StatusInternalServerError)
		}

		views.MakeIndex(races).Render(r.Context(), w)

	}
}

func handleUploadForm(w http.ResponseWriter, r *http.Request) {
	log.Printf("handleUploadForm ...")
	views.MakeUploadPage().Render(r.Context(), w)
}
