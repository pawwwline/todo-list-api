package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"todo-list-api/internal/logger"
	"todo-list-api/lib/e"

	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
	JwtSecret string
}

func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{
		JwtSecret: jwtSecret,
	}
}

func CreateToken(jwtSecret string, userId int) (string, error) {
	//token need to be refreshed after 48 hours
	if userId == 0 {
		return "", errors.New("user id is 0")
	}
	JwtPayload := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(48 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtPayload)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func GetTokenString(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		return "", e.AuthorizationMissing
	}
	if !strings.HasPrefix(reqToken, "Bearer ") {
		return "", e.AuthorizationMissing
	}
	reqToken = strings.TrimSpace(strings.TrimPrefix(reqToken, "Bearer "))
	return reqToken, nil
}

func ValidateToken(reqToken, jwtSecret string) (*jwt.Token, error) {
	logger.Logger.Debug("jwt secret for valid token", "secret", jwtSecret)
	logger.Logger.Debug("reqToken", "string", reqToken)
	token, err := jwt.Parse(reqToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil

}

func GetUserIdToken(token *jwt.Token) (int, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userId, exists := claims["sub"]; exists {
			if id, ok := userId.(float64); ok {
				logger.Logger.Debug("user id", "id", id)
				return int(id), nil
			}
			return 0, fmt.Errorf("error converting user_id to int")
		}
		return 0, fmt.Errorf("user_id not found in claims")
	}
	return 0, fmt.Errorf("invalid token")
}

func GetUserCtx(ctx context.Context) (int, error) {
	userId := ctx.Value("userId").(int)
	if userId == 0 {
		return -1, fmt.Errorf("empty id")
	}
	return userId, nil
}
