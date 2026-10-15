package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrNoAuthHeader = errors.New("no auth-header included in request")

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeader
	}

	authHeaderSplit := strings.Fields(authHeader)
	if len(authHeaderSplit) < 2 || strings.ToLower(authHeaderSplit[0]) != "bearer" {
		return "", fmt.Errorf("Bearer not found")
	}
	token := authHeaderSplit[1]

	return token, nil

}
