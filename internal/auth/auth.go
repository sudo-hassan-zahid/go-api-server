package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sudo-hassan-zahid/go-api-server/internal/config"
	"github.com/sudo-hassan-zahid/go-api-server/internal/database"
	customErr "github.com/sudo-hassan-zahid/go-api-server/internal/errors"
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

func GenerateRefreshToken(userID, role string) (string, error) {
	token, err := GenerateJWT(userID, role, RefreshTokenTTL)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	val, _ := database.Rdb.Get(context.Background(), "blacklist:"+tokenString).Result()
	if val != "" {
		return nil, customErr.ErrTokenInvalid
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, customErr.ErrTokenInvalid
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, customErr.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, customErr.ErrTokenInvalid
	}

	return claims, nil
}

func GenerateJWT(userID, role string, ttl time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func InvalidateToken(tokenString string, ttl time.Duration) error {
	return database.Rdb.Set(context.Background(), "blacklist:"+tokenString, "invalid", ttl).Err()
}
