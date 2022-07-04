package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Encrypt(password string) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt.GenerateFromPassword error: %v", err)
	}

	return string(encryptedPassword), nil
}

func Check(password string) error {
	hash, err := Encrypt(password)
	if err != nil {
		return fmt.Errorf("encrypt error: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("bcrypt.CompareHashAndPassword error: %v", err)
	}

	return nil
}
