package service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"pim-api-go/internal/middleware"
	"pim-api-go/internal/repository"
)

type AuthService struct{ repo *repository.AdminRepo }

func NewAuthService(repo *repository.AdminRepo) *AuthService { return &AuthService{repo} }

type LoginResult struct {
	Token    string
	Username string
}

func (s *AuthService) Login(username, password string) (*LoginResult, error) {
	admin, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("username atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("username atau password salah")
	}

	_ = s.repo.UpdateLastLogin(admin.ID)

	claims := &middleware.Claims{
		ID:       admin.ID,
		Username: admin.Username,
		Role:     string(admin.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	return &LoginResult{Token: signed, Username: admin.Username}, nil
}

func (s *AuthService) GetMe(id string) (string, string, error) {
	admin, err := s.repo.FindByID(id)
	if err != nil {
		return "", "", err
	}
	return admin.Username, string(admin.Role), nil
}
