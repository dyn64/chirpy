package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dyn64/chirpy/internal/auth"
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
		Body string `json:"body"`
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
		return
	}
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error in header: ", err)
		log.Printf("Getbearertoken failed: %v", err)
	}
	user_id, err := auth.ValidateJWT(bearer, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong token", err)
		log.Printf("Error validating jwt token, %v", err)
	}

	if len(input.Body) > 400 {
		respondWithError(w, http.StatusBadRequest, "Error, chirp is too long", nil)
		return
	}

	cleanResponse := responseClean{Clean: profanityFilter(input.Body)}
	createArgs := database.CreateChirpParams{
		Body:   cleanResponse.Clean,
		UserID: user_id,
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

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	arg := r.PathValue("ChirpID")
	if arg == "" {
		log.Printf("Error, missing chirp id in request")
		respondWithError(w, http.StatusBadRequest, "Missing ID in request", nil)
		return
	}
	id, err := uuid.Parse(arg)
	if err != nil {
		log.Printf("Error parsing ID: %v", err)
		respondWithError(w, http.StatusBadRequest, "Error parsing ID", err)
		return
	}
	chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), id)
	if err != nil {
		log.Printf("Error getting chirp: %s", err)
		errorMessage := fmt.Sprintf("Chirp with ID: %v not found", id)
		respondWithError(w, http.StatusNotFound, errorMessage, err)
	}

	actualChirp := chirpJSON{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, actualChirp)
}
