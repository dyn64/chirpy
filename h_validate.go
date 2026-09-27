package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerValidate(w http.ResponseWriter, r *http.Request) {
	type inputVals struct {
		Body string `json:"body"`
	}
	type responseOK struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	input := inputVals{}
	err := decoder.Decode(&input)
	if err != nil {
		log.Printf("Error reading validation input: %s", err)
		w.WriteHeader(500)
		return
	}
	if len(input.Body) > 400 {
		respondWithError(w, 400, "Error, chirp is too long")
		return
	}

	respondWithJSON(w, 200, responseOK{Valid: true})
}
