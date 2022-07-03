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
	hash, _ := Encrypt(password)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
