package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/dyn64/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func main() {
	godotenv.Load()
	const port = "8080"
	const rootpath = "."

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to db - %s", err)
	}

	mux := http.NewServeMux()

	fserv := http.FileServer(http.Dir(rootpath))
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      database.New(db),
		platform:       os.Getenv("PLATFORM"),
	}

	mux.Handle("/app/", apiCfg.middleWareMedtricsInc(http.StripPrefix("/app", fserv)))
	mux.HandleFunc("GET /api/healthz", handleReadyness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	//mux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidate)
	mux.HandleFunc("POST /api/chirps", apiCfg.addChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerAdduser)
	mux.HandleFunc("GET /api/chirps", apiCfg.getAllChirps)
	mux.HandleFunc("GET /api/chirps/{ChirpID}", apiCfg.getChirp)
	mux.HandleFunc("POST /api/login", apiCfg.handleLoginUser)

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Starting server on port: %s\n", port)
	log.Fatal(s.ListenAndServe())
}
