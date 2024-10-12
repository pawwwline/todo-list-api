package utils_test

import (
	"testing"
	"todo-list-api/internal/service/utils"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {

	tests := []struct {
		Name           string
		InputPassword  string
		ExpectError    bool
		ExpectedLength int
	}{{
		Name:           "Success Case",
		InputPassword:  "password123",
		ExpectError:    false,
		ExpectedLength: 60,
	},
		{
			Name:           "Empty Password",
			InputPassword:  "",
			ExpectError:    false,
			ExpectedLength: 60,
		}, {
			Name:           "Short Password",
			InputPassword:  "1",
			ExpectError:    false,
			ExpectedLength: 60,
		},
		{
			Name:           "Long Password",
			InputPassword:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis quis mi sem metus.",
			ExpectError:    true,
			ExpectedLength: 0,
		},
		{
			Name:           "Password length",
			InputPassword:  "Lorem ipsum dolor sit amet, consecteturt",
			ExpectError:    false,
			ExpectedLength: 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			hashedPassword, err := utils.HashPassword(tt.InputPassword)
			length := len(hashedPassword)
			if err == nil && tt.ExpectError {
				t.Errorf("expected error got nil")
			}
			if err != nil && !tt.ExpectError {
				t.Errorf("expected no error got %v", err)
			}
			if !tt.ExpectError {
				err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(tt.InputPassword))
				if err != nil {
					t.Errorf("expected password and hash to match")
				}
			}
			if length != tt.ExpectedLength {
				t.Errorf("expected length 60 got %v", length)
			}
		})
	}
}

func TestCompareHashPassword(t *testing.T) {
	password := "mySecretPassword"
	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	tests := []struct {
		Name          string
		InputPassword string
		InputHash     string
		ExpectError   bool
	}{
		{
			Name:          "Success case",
			InputPassword: password,
			InputHash:     hash,
			ExpectError:   false,
		},
		{
			Name:          "Not matching",
			InputPassword: password + "1",
			InputHash:     hash,
			ExpectError:   true,
		},
		{
			Name:          "Empty hash",
			InputPassword: password,
			InputHash:     "",
			ExpectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := utils.CompareHashPassword(tt.InputPassword, tt.InputHash)
			if err == nil && tt.ExpectError {
				t.Errorf("expected error got nil")
			}
			if err != nil && !tt.ExpectError {
				t.Errorf("expected no error got %v, %s", err, tt.InputHash)
			}
		})
	}
}
