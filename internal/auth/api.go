package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeader
	}

	authHeaderSplit := strings.Fields(authHeader)
	if len(authHeaderSplit) < 2 || strings.ToLower(authHeaderSplit[0]) != "apikey" {
		return "", fmt.Errorf("ApiKey not found")
	}
	token := authHeaderSplit[1]

	return token, nil

}
