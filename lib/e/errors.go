package e

import (
	"errors"
	"fmt"
)

var (
	UniqueViolationErr    = errors.New("value already exist in db")
	AuthorizationMissing  = errors.New("authorization header is missing")
	InvalidCredentialsErr = errors.New("invalid email or password")
)

func WrapError(message string, err error) error {
	if err != nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}
