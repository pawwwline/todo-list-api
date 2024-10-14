package e

import (
	"errors"
	"fmt"
)

var (
	UniqueViolationErr    = errors.New("value already exist in db")
	AuthorizationMissing  = errors.New("authorization header is missing")
	InvalidCredentialsErr = errors.New("invalid email or password")
	ItemIdNotFound        = errors.New("item id not found")
)

func WrapError(message string, err error) error {
	if err != nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}
