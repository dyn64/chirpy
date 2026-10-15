package main

import (
	"net/http"

	"github.com/dyn64/chirpy/internal/auth"
)

func (cfg *apiConfig) RefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	refreshtoken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error in header, missing token", err)
		return
	}
	user, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refreshtoken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized token", err)
		return
	}
	// 	if full_token.RevokedAt.Valid {
	// 		respondWithError(w, http.StatusUnauthorized, "Token revoked", err)
	// 		return
	// 	}
	// 	if full_token.ExpiresAt.Before(time.Now()) {
	// 		respondWithError(w, http.StatusUnauthorized, "Token expired", err)
	// 		return
	// 	}
	//
	new_token, err := auth.MakeJWT(user.ID, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error making new jwt: %v", err)
		return
	}

	resp_token := response{
		Token: new_token,
	}
	respondWithJSON(w, http.StatusOK, resp_token)
}

func (cfg *apiConfig) RevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshtoken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error in header, missing token", err)
		return
	}

	// full_token, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refreshtoken)
	// if err != nil || full_token.RevokedAt.Valid || full_token.ExpiresAt.Before(time.Now()) {
	// 	respondWithError(w, http.StatusUnauthorized, "Unauthorized token", err)
	// 	return
	// }

	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), refreshtoken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking token: %v", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, "")
}
