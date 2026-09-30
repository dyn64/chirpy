package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dyn64/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) addChirp(w http.ResponseWriter, r *http.Request) {
	type inputVals struct {
		Body    string    `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}
	type responseOK struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
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
		Body: cleanResponse.Clean,
		UserID: uuid.NullUUID{
			UUID:  input.User_id,
			Valid: true,
		},
	}
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), createArgs)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error adding chirp to the database: %s", err)
		return
	}
	actualChirp := responseOK{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID.UUID,
	}
	respondWithJSON(w, http.StatusCreated, actualChirp)

	// respondWithJSON(w, http.StatusOK, cleanResponse)
}
