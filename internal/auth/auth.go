package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	appError "github.com/sudo-hassan-zahid/go-api-server/internal/errors"
)

type JWT struct {
	secret []byte
	ttl    time.Duration
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func New(secret []byte, ttl time.Duration) *JWT {
	return &JWT{
		secret: secret,
		ttl:    ttl,
	}
}

func (j *JWT) Generate(userID, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWT) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, appError.ErrTokenInvalid
			}
			return j.secret, nil
		},
	)

	if err != nil || !token.Valid {
		return nil, appError.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, appError.ErrTokenInvalid
	}

	return claims, nil
}
