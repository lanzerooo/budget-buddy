package auth

import (
	"budgetbuddy/pkg/config"
	"budgetbuddy/pkg/logger"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(email string) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config: ", err)
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(cfg.JWTSecret))
}
