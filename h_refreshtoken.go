package main

import (
	"net/http"
	"time"

	"github.com/dyn64/chirpy/internal/auth"
)

func (cfg *apiConfig) RefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	refreshtoken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error in header, missing token", err)
	}
	full_token, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refreshtoken)
	if err != nil || full_token.RevokedAt.Valid || full_token.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized token", err)
	}

	new_token, err := auth.MakeJWT(full_token.UserID, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error making new jwt: %v", err)
	}

	respondWithJSON(w, http.StatusOK, new_token)
}

func (cfg *apiConfig) RevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshtoken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error in header, missing token", err)
	}

	full_token, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refreshtoken)
	if err != nil || full_token.RevokedAt.Valid || full_token.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized token", err)
	}

	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), full_token.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking token: %v", err)
	}
	respondWithJSON(w, http.StatusNoContent, "")
}
