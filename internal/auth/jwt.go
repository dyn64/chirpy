package auth

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claimsStruct := &jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsStruct)
	signed_token, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signed_token, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	type claimStruct struct {
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(tokenString, &claimStruct{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("Failed to get subject from claim: %v", err)
		return uuid.Nil, nil
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}
	ID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Invalid user ID: %w", err)
	}

	return ID, nil
}
