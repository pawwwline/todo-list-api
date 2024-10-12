package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
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
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) != 2 {
		return "", e.AuthorizationMissing
	}
	reqToken = strings.TrimSpace(splitToken[1])
	return reqToken, nil
}

func ValidateToken(jwtSecret, reqToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(reqToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil

}

func GetUserIdToken(token *jwt.Token) (int, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["sub"].(int)
		if !ok {
			return -1, fmt.Errorf("error getting user_id from jwt")
		}
		return userId, nil
	}
	return -1, fmt.Errorf("error getting user_id from jwt")

}

func GetUserCtx(ctx context.Context) (int, error) {
	userId := ctx.Value("userId").(int)
	if userId == 0 {
		return -1, fmt.Errorf("empty id")
	}
	return userId, nil
}
