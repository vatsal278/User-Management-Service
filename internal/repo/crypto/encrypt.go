package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password []byte, salt []byte) (string, error) {
	password = append(password, salt...)
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
