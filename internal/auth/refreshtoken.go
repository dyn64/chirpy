package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	r_num := make([]byte, 32)
	rand.Read(r_num)
	return hex.EncodeToString(r_num)
}
