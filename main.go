package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	const rootpath = "."

	mux := http.NewServeMux()

	fserv := http.FileServer(http.Dir(rootpath))
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux.Handle("/app/", apiCfg.middleWareMedtricsInc(http.StripPrefix("/app", fserv)))
	mux.HandleFunc("GET /api/healthz", handleReadyness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.middleWareMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

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
