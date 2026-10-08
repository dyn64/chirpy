package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	authHeaderSplit := strings.Fields(authHeader)
	if len(authHeaderSplit) != 2 || strings.ToLower(authHeaderSplit[0]) != "bearer" {
		return "", fmt.Errorf("Bearer not found")
	}
	token := authHeaderSplit[1]

	return token, nil

}
