package auth

import (
	"context"

	"github.com/sudo-hassan-zahid/go-api-server/internal/database"
	"github.com/sudo-hassan-zahid/go-api-server/internal/errors"
	"github.com/sudo-hassan-zahid/go-api-server/internal/models"
	"github.com/sudo-hassan-zahid/go-api-server/internal/repository"
	"github.com/sudo-hassan-zahid/go-api-server/utils"
)

type Service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) Register(email, password, firstName, lastName string) (*models.User, error) {
	existingUser, _ := s.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	user := &models.User{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(email, password string) (string, string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", "", errors.ErrInvalidCredentials
	}

	if !utils.CheckPassword(user.Password, password) {
		return "", "", errors.ErrInvalidCredentials
	}

	if err := s.handleMaxLogins(user.ID.String()); err != nil {
		return "", "", err
	}

	accessToken, err := GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken(user.ID.String(), user.Role)
	if err != nil {
		return "", "", err
	}

	if err := s.storeSession(user.ID.String(), refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) handleMaxLogins(userID string) error {
	ctx := context.Background()
	key := "sessions:" + userID

	count, err := database.Rdb.LLen(ctx, key).Result()
	if err != nil {
		return err
	}

	if count >= 5 {
		oldestToken, err := database.Rdb.LPop(ctx, key).Result()
		if err != nil {
			return err
		}
		InvalidateToken(oldestToken, RefreshTokenTTL)
	}

	return nil
}

func (s *Service) storeSession(userID, refreshToken string) error {
	key := "sessions:" + userID
	return database.Rdb.RPush(context.Background(), key, refreshToken).Err()
}

func (s *Service) Logout(userID, tokenString string) error {
	if err := InvalidateToken(tokenString, AccessTokenTTL); err != nil {
		return err
	}

	return nil
}

func (s *Service) LogoutAll(userID string) error {
	ctx := context.Background()
	key := "sessions:" + userID

	tokens, err := database.Rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return err
	}

	for _, t := range tokens {
		InvalidateToken(t, RefreshTokenTTL)
	}

	return database.Rdb.Del(ctx, key).Err()
}

func (s *Service) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := ValidateToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	ctx := context.Background()
	key := "sessions:" + claims.UserID

	tokens, err := database.Rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return "", "", err
	}
	found := false
	for _, t := range tokens {
		if t == refreshToken {
			found = true
			break
		}
	}

	if !found {
		return "", "", errors.ErrTokenInvalid
	}

	newAccessToken, err := GenerateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := GenerateRefreshToken(claims.UserID, claims.Role)
	if err != nil {
		return "", "", err
	}

	if err := database.Rdb.LRem(ctx, key, 1, refreshToken).Err(); err != nil {
		return "", "", err
	}
	if err := database.Rdb.RPush(ctx, key, newRefreshToken).Err(); err != nil {
		return "", "", err
	}

	InvalidateToken(refreshToken, RefreshTokenTTL)

	return newAccessToken, newRefreshToken, nil
}
