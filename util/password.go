package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GetHashedPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("Generate password failed: %w", err)
	}
	return string(hashed), nil
}

func CheckPassword(password string, hashed []byte) error {
	return bcrypt.CompareHashAndPassword(hashed, []byte(password))
}
