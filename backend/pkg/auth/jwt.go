package auth

import (
	"time"

	"budgetbuddy/pkg/config"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(email string) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}
