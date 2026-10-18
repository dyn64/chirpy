package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) middleWareMedtricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	count := fmt.Sprintf("<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>", cfg.fileserverHits.Load())
	w.Write([]byte(count))
}

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type webhook struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	input := webhook{}
	err := decoder.Decode(&input)
	if err != nil || input.Event != "user.upgraded" {
		log.Printf("Error in webhook or wrong event")
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	user, err := cfg.dbQueries.UpgradeUser(r.Context(), input.Data.UserID)
	if err != nil || user.ID != input.Data.UserID {
		respondWithError(w, http.StatusNotFound, "", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
