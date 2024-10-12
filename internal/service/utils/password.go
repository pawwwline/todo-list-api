package utils

import "golang.org/x/crypto/bcrypt"

type HashingPasswords interface {
	HashPassword(password string) (string, error)
	CompareHashPassword(password, hash string) error
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CompareHashPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
