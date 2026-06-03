package repository

import (
	"pim-api-go/internal/model"

	"gorm.io/gorm"
)

type AdminRepo struct{ db *gorm.DB }

func NewAdminRepo(db *gorm.DB) *AdminRepo { return &AdminRepo{db} }

func (r *AdminRepo) FindByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	if err := r.db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepo) FindByID(id string) (*model.Admin, error) {
	var admin model.Admin
	if err := r.db.First(&admin, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepo) UpdateLastLogin(id string) error {
	return r.db.Model(&model.Admin{}).Where("id = ?", id).
		Update("lastLogin", gorm.Expr("NOW()")).Error
}

func (r *AdminRepo) FindAll() ([]model.Admin, error) {
	var admins []model.Admin
	if err := r.db.Order(`"createdAt" asc`).Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}

func (r *AdminRepo) Create(admin *model.Admin) error {
	return r.db.Create(admin).Error
}

func (r *AdminRepo) UpdateRole(id string, role model.AdminRole) error {
	return r.db.Model(&model.Admin{}).Where("id = ?", id).Update("role", role).Error
}

func (r *AdminRepo) UpdatePassword(id, hash string) error {
	return r.db.Model(&model.Admin{}).Where("id = ?", id).Update("password_hash", hash).Error
}

func (r *AdminRepo) Delete(id string) error {
	return r.db.Delete(&model.Admin{}, "id = ?", id).Error
}
