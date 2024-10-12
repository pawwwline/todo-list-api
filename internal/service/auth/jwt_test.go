package auth_test

import (
	"testing"
	"todo-list-api/internal/service/auth"
)

func TestCreateJWT(t *testing.T) {
	secret := "secret"

	token, err := auth.CreateToken(secret, 1)
	if err != nil {
		t.Errorf("error creating JWT: %v", err)
	}

	if token == "" {
		t.Error("expected token to be not empty")
	}
}
