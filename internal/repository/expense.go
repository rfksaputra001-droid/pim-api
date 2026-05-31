package repository

import (
	"time"

	"pim-api-go/internal/model"

	"gorm.io/gorm"
)

type ExpenseRepo struct{ db *gorm.DB }

func NewExpenseRepo(db *gorm.DB) *ExpenseRepo { return &ExpenseRepo{db} }

func (r *ExpenseRepo) FindPublic(page, limit int) ([]model.Expense, int64, error) {
	var total int64
	r.db.Model(&model.Expense{}).Count(&total)

	var items []model.Expense
	err := r.db.Order("date desc").Offset((page - 1) * limit).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *ExpenseRepo) FindAdmin(page, limit int) ([]model.Expense, int64, error) {
	var total int64
	r.db.Model(&model.Expense{}).Count(&total)

	var items []model.Expense
	err := r.db.Order("date desc").Offset((page - 1) * limit).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *ExpenseRepo) FindByID(id string) (*model.Expense, error) {
	var e model.Expense
	if err := r.db.First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *ExpenseRepo) Create(e *model.Expense) error {
	return r.db.Create(e).Error
}

func (r *ExpenseRepo) SumAll() (int64, error) {
	var result struct{ Sum int64 }
	err := r.db.Model(&model.Expense{}).
		Select("COALESCE(SUM(amount), 0) as sum").Scan(&result).Error
	return result.Sum, err
}

func (r *ExpenseRepo) FindChartData(since time.Time) ([]model.Expense, error) {
	var items []model.Expense
	err := r.db.Where("date >= ?", since).Select("amount, date").Find(&items).Error
	return items, err
}
