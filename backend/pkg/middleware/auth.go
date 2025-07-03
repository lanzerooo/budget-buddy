package middleware

import (
	"budgetbuddy/pkg/logger"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func AuthMiddleware(jwtSecret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			logger.Error("Authorization header is empty")
			return
		}
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		} else {
			http.Error(w, "Authorization header must start with 'Bearer '", http.StatusUnauthorized)
			logger.Error("Invalid Authorization header format")
			return
		}

		if tokenStr == "" {
			http.Error(w, "JWT token is empty", http.StatusUnauthorized)
			logger.Error("JWT token is empty after removing Bearer prefix")
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			logger.Error("Failed to parse or validate token: ", err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			logger.Error("Failed to parse token claims")
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			http.Error(w, "Invalid email in token", http.StatusUnauthorized)
			logger.Error("Email not found in token claims")
			return
		}

		r.Header.Set("X-User-Email", email)
		next(w, r)
	}
}
