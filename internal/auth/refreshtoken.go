package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	var r_num []byte
	rand.Read(r_num)
	return hex.EncodeToString(r_num)
}
