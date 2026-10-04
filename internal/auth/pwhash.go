package auth

import (
	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	argonparams := argon2id.DefaultParams
	hash, err := argon2id.CreateHash(password, argonparams)
	if err != nil {
		//log.Printf("Error creating password-hash: %v", err)
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	pwHash, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		//log.Printf("Error checking password hash: %v", err)
		return false, err
	}
	return pwHash, nil

}
