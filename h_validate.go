package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func (cfg *apiConfig) handlerValidate(w http.ResponseWriter, r *http.Request) {
	type inputVals struct {
		Body string `json:"body"`
	}
	type responseOK struct {
		Valid bool `json:"valid"`
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
		w.WriteHeader(500)
		return
	}
	if len(input.Body) > 400 {
		respondWithError(w, http.StatusBadRequest, "Error, chirp is too long", nil)
		return
	}

	cleanResponse := responseClean{Clean: profanityFilter(input.Body)}

	respondWithJSON(w, 200, cleanResponse)
}

func profanityFilter(s string) string {
	filteredWords := regexp.MustCompile(`kerfuffle|sharbert|fornax`)

	words := strings.Split(s, " ")
	for i, word := range words {
		if filteredWords.MatchString(strings.ToLower(word)) {
			words[i] = "****"
		}

	}
	return strings.Join(words, " ")
}
