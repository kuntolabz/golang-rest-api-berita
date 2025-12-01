package services

import (
	"errors"

	"github.com/kunto/golang-rest-api-berita/repositories"
	"github.com/kunto/golang-rest-api-berita/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authService{userRepo: repo}
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Cek password hash
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate token
	token, err := utils.GenerateToken(user.IdUser, user.Email)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
