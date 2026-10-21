package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
		return
	}
	user_id, err := auth.ValidateJWT(bearer, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong token", err)
		log.Printf("Error validating jwt token, %v", err)
		return
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
	auth := r.URL.Query().Get("author_id")
	var chirps []database.Chirp
	var err error
	var Sorting string

	sort_par := r.URL.Query().Get("sort")
	if sort_par == "" || strings.ToLower(sort_par) == "asc" {
		Sorting = "ASC"

	} else if strings.ToLower(sort_par) == "desc" {
		Sorting = "DESC"
	}
	if auth != "" {
		author_id, err := uuid.Parse(auth)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Error in author-id", err)
			return
		}
		chirps, err = cfg.dbQueries.GetChirpsByAuthorID(r.Context(), database.GetChirpsByAuthorIDParams{UserID: author_id, Sorting: Sorting})
	} else {
		chirps, err = cfg.dbQueries.GetChirps(r.Context(), Sorting)
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps: %v", err)
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
		return
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

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	arg := r.PathValue("chirpID")
	chirp_id, err := uuid.Parse(arg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing chirp-id", err)
		return
	}
	// if arg == "" {
	// 	log.Printf("Error, missing chirp id in (delete) request")
	// 	respondWithError(w, http.StatusBadRequest, "Missing ID in request", nil)
	// 	return
	// }

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error, missing access token:", err)
		return
	}
	user_id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chirp_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	if chirp.UserID != user_id {
		respondWithError(w, http.StatusForbidden, "Verboten", err)
		return
	}

	err = cfg.dbQueries.DeleteChirp(r.Context(), chirp_id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't delete chirp", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
