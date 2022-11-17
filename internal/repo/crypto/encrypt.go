package crypto

import (
	"crypto/sha512"
	"encoding/hex"
)

func GeneratePasswordHash(password []byte, salt []byte) (string, error) {
	password = append(password, salt...)
	var sha512Hasher = sha512.New()
	sha512Hasher.Write(password)
	var hashedPasswordBytes = sha512Hasher.Sum(nil)
	// Convert the hashed password to a hex string
	hash := hex.EncodeToString(hashedPasswordBytes)

	return hash, nil
}
