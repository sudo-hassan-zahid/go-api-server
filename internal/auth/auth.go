package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sudo-hassan-zahid/go-api-server/internal/config"
	appError "github.com/sudo-hassan-zahid/go-api-server/internal/errors"
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

var jwtSecret []byte

func Init(cfg *config.Config) {
	if len(cfg.App.JWTSecret) == 0 {
		panic("JWT_SECRET is required")
	}
	jwtSecret = cfg.App.JWTSecret
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID, role string) (string, error) {
	return GenerateJWT(userID, role, AccessTokenTTL)
}

func GenerateRefreshToken(userID string) (string, error) {
	return GenerateJWT(userID, "", RefreshTokenTTL)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, appError.ErrTokenInvalid
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, appError.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, appError.ErrTokenInvalid
	}

	return claims, nil
}

func GenerateJWT(userID, role string, ttl time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
