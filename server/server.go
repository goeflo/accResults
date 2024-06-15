package server

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/goeflo/accResults/config"
	"github.com/goeflo/accResults/internal/data"
	"github.com/goeflo/accResults/internal/database"
	"github.com/goeflo/accResults/server/views"
	"github.com/gorilla/mux"
)

type Server struct {
	DB     database.Database
	Config *config.RaceResultConfig
	Router *mux.Router
}

const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB

func NewServer(config *config.RaceResultConfig, db database.Database) Server {

	router := mux.NewRouter()

	router.HandleFunc("/", handleHome).Methods("GET")
	router.HandleFunc("/upload", handleUploadForm).Methods("GET")
	router.HandleFunc("/upload", handleUpload(config, db)).Methods("POST")

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
	slog.Info("starting http server", "addr", s.Config.HttpServerAddr)
	http.ListenAndServe(s.Config.HttpServerAddr, s.Router)
}

func handleUpload(c *config.RaceResultConfig, db database.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("handleUpload ...")

		r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
		if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
			slog.Error(err.Error())
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
			slog.Error(err.Error())
			http.Error(w, "result file was not saved "+err.Error(), http.StatusInternalServerError)
		}

		db.NewResult(result.Filename)
		//fmt.Printf("result file %v\n", string(result))

		// TODO show upload result
		views.MakeIndex().Render(r.Context(), w)

	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	slog.Info("handleHome ...")
	views.MakeIndex().Render(r.Context(), w)
}

func handleUploadForm(w http.ResponseWriter, r *http.Request) {
	slog.Info("handleUploadForm ...")
	views.MakeUploadPage().Render(r.Context(), w)
}
