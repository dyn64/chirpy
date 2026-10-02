package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dyn64/chirpy/internal/database"
	"github.com/google/uuid"
)

type chirpJSON struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) addChirp(w http.ResponseWriter, r *http.Request) {
	type inputVals struct {
		Body    string    `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}
	type responseClean struct {
		Clean string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	input := inputVals{}
	err := decoder.Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		log.Printf("Error reading validation input: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(input.Body) > 400 {
		respondWithError(w, http.StatusBadRequest, "Error, chirp is too long", nil)
		return
	}

	cleanResponse := responseClean{Clean: profanityFilter(input.Body)}
	createArgs := database.CreateChirpParams{
		Body:   cleanResponse.Clean,
		UserID: input.User_id,
	}

	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), createArgs)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error adding chirp to the database: %s", err)
		return
	}
	actualChirp := chirpJSON{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, actualChirp)

}

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps: %s", err)
		return
	}
	var aChirps []chirpJSON
	for _, chirp := range chirps {
		aChirps = append(aChirps, chirpJSON{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, aChirps)
}
