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
	var apiCfg apiConfig

	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middleWareMedtricsInc(fserv)))
	mux.HandleFunc("/healthz", handleReadyness)
	mux.Handle("/metrics", http.StripPrefix("/metrics", apiCfg.middleWareResponse(fserv)))
	mux.Handle("/reset", http.StripPrefix("/reset", apiCfg.middleWareReset(fserv)))

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
