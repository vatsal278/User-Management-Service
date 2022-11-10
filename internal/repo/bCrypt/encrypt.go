package bCrypt

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// GeneratePasswordHash handles generating password hash
// bCrypt hashes password of type byte
func GeneratePasswordHash(password []byte) string {
	// default cost is 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	// return stringified password
	return string(hashedPassword)
}

// PasswordCompare handles password hash compare
func PasswordCompare(password []byte, hashedPassword []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)

	return err
}
