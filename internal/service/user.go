package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"pim-api-go/internal/model"
	"pim-api-go/internal/repository"
)

type UserService struct{ repo *repository.AdminRepo }

func NewUserService(repo *repository.AdminRepo) *UserService { return &UserService{repo} }

func (s *UserService) List() ([]model.Admin, error) {
	return s.repo.FindAll()
}

func (s *UserService) Create(username, password string, role model.AdminRole) (*model.Admin, error) {
	if strings.TrimSpace(username) == "" {
		return nil, fmt.Errorf("Username wajib diisi")
	}
	if len(password) < 8 {
		return nil, fmt.Errorf("Password minimal 8 karakter")
	}
	if role != model.RoleAdmin && role != model.RoleSuperAdmin {
		return nil, fmt.Errorf("Role tidak valid")
	}

	existing, _ := s.repo.FindByUsername(username)
	if existing != nil {
		return nil, fmt.Errorf("Username sudah digunakan")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	admin := &model.Admin{
		ID:           uuid.NewString(),
		Username:     strings.TrimSpace(username),
		PasswordHash: string(hash),
		Role:         role,
		CreatedAt:    time.Now(),
	}
	if err := s.repo.Create(admin); err != nil {
		return nil, err
	}
	return admin, nil
}

func (s *UserService) UpdateRole(id, requestorID string, role model.AdminRole) error {
	if role != model.RoleAdmin && role != model.RoleSuperAdmin {
		return fmt.Errorf("Role tidak valid")
	}
	if id == requestorID {
		return fmt.Errorf("Tidak bisa mengubah role sendiri")
	}
	return s.repo.UpdateRole(id, role)
}

func (s *UserService) UpdatePassword(id, newPassword string) error {
	if len(newPassword) < 8 {
		return fmt.Errorf("Password minimal 8 karakter")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(id, string(hash))
}

func (s *UserService) Delete(id, requestorID string) error {
	if id == requestorID {
		return fmt.Errorf("Tidak bisa menghapus akun sendiri")
	}
	target, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("Pengguna tidak ditemukan")
	}
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	_ = target
	return nil
}
