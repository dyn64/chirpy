package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) middleWareReset(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Store(0)
	})
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset only allowed in dev environment"))
		return
	}
	err := cfg.dbQueries.ResetUserTable(r.Context())
	if err != nil {
		log.Printf("Error resetting users-table: %s", err)
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database reset and hitcounter set to 0."))
}
