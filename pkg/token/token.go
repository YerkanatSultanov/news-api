package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	Secret          string
	ExpirationHours int
}

func NewJWTManager(secret string, expirationHours int) *JWTManager {
	return &JWTManager{Secret: secret, ExpirationHours: expirationHours}
}

func (j *JWTManager) GenerateTokens(userID int, role string) (accessToken string, refreshToken string, err error) {
	now := time.Now()

	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     now.Add(time.Duration(j.ExpirationHours) * time.Hour).Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = access.SignedString([]byte(j.Secret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     now.Add(7 * 24 * time.Hour).Unix(),
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refresh.SignedString([]byte(j.Secret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
