package middleware

import (
	"context"
	"net/http"
	"todo-list-api/internal/service/auth"
)

type contextKey string

const UserIdKey contextKey = "userID"

func TokenAuth(next http.Handler, jwtSecret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/register" && r.URL.Path != "/login" {
			tokenStr, err := auth.GetTokenString(r)
			if err != nil {
				http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
				return
			}
			tokenJWT, err := auth.ValidateToken(tokenStr, jwtSecret)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}
			userId, err := auth.GetUserIdToken(tokenJWT)
			if err != nil {
				http.Error(w, "Unable to extract user ID from token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIdKey, userId)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})

}

func TokenAuthMiddleware(secret string) Middleware {
	return func(next http.Handler) http.Handler {
		return TokenAuth(next, secret)
	}
}
