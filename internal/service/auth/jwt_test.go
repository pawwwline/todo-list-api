package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"testing"
	"time"
	"todo-list-api/internal/service/auth"
	"todo-list-api/lib/e"

	"github.com/dgrijalva/jwt-go"
)

func TestCreateJWT(t *testing.T) {
	secret := "secret"

	token, err := auth.CreateToken(secret, 1)
	if err != nil {
		t.Errorf("error creating JWT: %v", err)
	}
	t.Logf("token %v", token)

	if token == "" {
		t.Error("expected token to be not empty")
	}

}

func TestGetTokenString(t *testing.T) {
	// Тест: правильный токен
	reqWithToken := &http.Request{
		Header: http.Header{
			"Authorization": []string{"Bearer test_token_string"},
		},
	}

	token, err := auth.GetTokenString(reqWithToken)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if token != "test_token_string" {
		t.Errorf("expected 'test_token_string', got %v", token)
	}

	reqWithoutAuth := &http.Request{
		Header: http.Header{},
	}

	_, err = auth.GetTokenString(reqWithoutAuth)
	if err == nil {
		t.Error("expected error for missing Authorization header, got none")
	}
	if err != e.AuthorizationMissing {
		t.Errorf("expected error AuthorizationMissing, got %v", err)
	}

	reqWithWrongFormat := &http.Request{
		Header: http.Header{
			"Authorization": []string{"InvalidBearer test_token_string"},
		},
	}

	_, err = auth.GetTokenString(reqWithWrongFormat)
	if err == nil {
		t.Error("expected error for incorrect Authorization format, got none")
	}
	if err != e.AuthorizationMissing {
		t.Errorf("expected error AuthorizationMissing for incorrect format, got %v", err)
	}
}

func createRSAPrivateKey() *rsa.PrivateKey {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	return privateKey
}

func TestValidateToken(t *testing.T) {
	jwtSecret := "mysecret"

	createValidToken := func(userId int) string {
		token, err := auth.CreateToken(jwtSecret, userId)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}
		return token
	}

	t.Run("valid token", func(t *testing.T) {
		validToken := createValidToken(1)
		token, err := auth.ValidateToken(jwtSecret, validToken)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if token == nil || !token.Valid {
			t.Error("expected valid token, got invalid")
		}
	})

	t.Run("invalid signature", func(t *testing.T) {
		invalidSecret := "wrongsecret"
		validToken := createValidToken(1)
		_, err := auth.ValidateToken(invalidSecret, validToken)
		if err == nil {
			t.Error("expected error due to invalid signature, got none")
		}
	})

	t.Run("unexpected signing method", func(t *testing.T) {
		privateKey := createRSAPrivateKey()
		claims := jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 1).Unix(),
			"sub": "1234567890",
		}
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		tokenString, err := token.SignedString(privateKey)

		if err != nil {
			t.Fatalf("failed to sign token: %v", err)
		}

		_, err = auth.ValidateToken(jwtSecret, tokenString)
		if err == nil || err.Error() != "unexpected signing method: RS256" {
			t.Errorf("expected error for unexpected signing method, got %v", err)
		}
	})

	t.Run("malformed token", func(t *testing.T) {
		malformedToken := "thisisnotatoken"
		_, err := auth.ValidateToken(jwtSecret, malformedToken)
		if err == nil {
			t.Error("expected error for malformed token, got none")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		claims := jwt.MapClaims{
			"exp": time.Now().Add(-time.Hour).Unix(),
			"sub": "1234567890",
		}
		expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := expiredToken.SignedString([]byte(jwtSecret))

		_, err := auth.ValidateToken(jwtSecret, tokenString)
		if err == nil {
			t.Error("expected error for expired token, got none")
		}
	})
}
