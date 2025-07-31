package service

import (
	"errors"
	"mini-crm/internal/dto"
	"mini-crm/internal/models"
	"mini-crm/internal/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo      repository.UserRepository
	TokenRepo repository.TokenRepository
}

func NewUserService(repo repository.UserRepository, tokenRepo repository.TokenRepository) *UserService {
	return &UserService{
		Repo:      repo,
		TokenRepo: tokenRepo,
	}
}

func (s *UserService) Register(input dto.RegisterUserInput) error {
	existing, _ := s.Repo.GetByEmail(input.Email)
	if existing != nil {
		return errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	return s.Repo.Create(user)
}

func (s *UserService) Login(input dto.LoginUserInput) (map[string]string, error) {
	user, err := s.Repo.GetByEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, refreshToken, expiresAt, err := GenerateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	rt := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	}
	if err := s.TokenRepo.Save(rt); err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func GenerateTokens(userID uint) (accessToken string, refreshToken string, refreshExpiresAt time.Time, err error) {
	accessClaims := dto.CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = access.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return
	}

	refreshToken = uuid.NewString()
	refreshExpiresAt = time.Now().Add(7 * 24 * time.Hour)

	return
}

func (s *UserService) RefreshTokens(refreshToken string) (map[string]string, error) {
	tokenModel, err := s.TokenRepo.Find(refreshToken)
	if err != nil || tokenModel.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("invalid or expired refresh token")
	}

	accessToken, newRefresh, expiresAt, err := GenerateTokens(tokenModel.UserID)
	if err != nil {
		return nil, err
	}

	_ = s.TokenRepo.Delete(tokenModel.Token)
	_ = s.TokenRepo.Save(&models.RefreshToken{
		UserID:    tokenModel.UserID,
		Token:     newRefresh,
		ExpiresAt: expiresAt,
	})

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": newRefresh,
	}, nil
}

func (s *UserService) Logout(refreshToken string) error {
	return s.TokenRepo.Delete(refreshToken)
}

func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	return s.Repo.GetByID(userID)
}

func (s *UserService) GenerateTokens(userID uint) (string, string, time.Time, error) {
	return GenerateTokens(userID)
}
