package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func NormalizedPassword(p string) []byte {
	// Normalize the password
	return []byte(p)
}

func HashPassword(p string) string {
	// Hash the password
	pwd := NormalizedPassword(p)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

