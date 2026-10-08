package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dyn64/chirpy/internal/auth"
	"github.com/dyn64/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerAdduser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	input := request{}
	err := decoder.Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		log.Printf("Error reading email input: %s", err)
		return
	}
	pwhash, err := auth.HashPassword(input.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
	}
	dbCreateuserparams := database.CreateUserParams{
		Email:          input.Email,
		HashedPassword: pwhash,
	}
	usr, err := cfg.dbQueries.CreateUser(r.Context(), dbCreateuserparams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error adding user to database: %s", err)
		return
	}
	user := response{
		ID:        usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email:     usr.Email,
	}
	respondWithJSON(w, http.StatusCreated, user)

}

func (cfg *apiConfig) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Password string `json:"password"`
		Email    string `json:"email"`
		Expires  *int   `json:"expires_in_seconds,omitempty"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	input := request{}
	err := decoder.Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		log.Printf("Error reading email input: %s", err)
		return
	}

	usr, err := cfg.dbQueries.GetUserFromEmail(r.Context(), input.Email)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		respondWithError(w, http.StatusUnauthorized, "User not found", err)
		return
	}

	ok, err := auth.CheckPasswordHash(input.Password, usr.HashedPassword)
	if err != nil || !ok {
		log.Printf("Error in checkpassword: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or passord", err)
		return
	}
	var token_life time.Duration
	if input.Expires != nil {
		token_life = time.Second * time.Duration(*input.Expires)
	} else {
		token_life = time.Hour * 1
	}
	token, err := auth.MakeJWT(usr.ID, cfg.secret, token_life)

	user := response{
		ID:        usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email:     usr.Email,
		Token:     token,
	}

	respondWithJSON(w, http.StatusOK, user)

}
