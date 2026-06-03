package repository

import (
	"time"

	"pim-api-go/internal/model"

	"gorm.io/gorm"
)

type ContributionRepo struct{ db *gorm.DB }

func NewContributionRepo(db *gorm.DB) *ContributionRepo { return &ContributionRepo{db} }

type ContributionFilter struct {
	Status  model.ContributionStatus
	Page    int
	Limit   int
	SortBy  string
	SortDir string
}

func (r *ContributionRepo) FindPublic(year, month, page, limit int) ([]model.Contribution, int64, error) {
	query := r.db.Model(&model.Contribution{}).Where("status = ?", model.StatusVerified)

	if year > 0 && month >= 1 && month <= 12 {
		start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, 0)
		query = query.Where("verified_at >= ? AND verified_at < ?", start, end)
	} else if year > 0 {
		start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(1, 0, 0)
		query = query.Where("verified_at >= ? AND verified_at < ?", start, end)
	}

	var total int64
	query.Count(&total)

	var items []model.Contribution
	offset := (page - 1) * limit
	err := query.Select(`id, name, amount, notes, "createdAt", "verifiedAt"`).
		Order(`"verifiedAt" desc`).Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *ContributionRepo) FindAdmin(f ContributionFilter) ([]model.Contribution, int64, error) {
	query := r.db.Model(&model.Contribution{})
	if f.Status != "" {
		query = query.Where("status = ?", f.Status)
	}

	// frontend kirim camelCase, map ke kolom DB yang benar
	sortMap := map[string]string{
		"createdAt": `"createdAt"`,
		"amount":    "amount",
		"name":      "name",
	}
	sortBy := `"createdAt"`
	if col, ok := sortMap[f.SortBy]; ok {
		sortBy = col
	}
	sortDir := "desc"
	if f.SortDir == "asc" {
		sortDir = "asc"
	}

	var total int64
	query.Count(&total)

	p := max(f.Page, 1)
	var items []model.Contribution
	offset := (p - 1) * f.Limit
	err := query.Order(sortBy + " " + sortDir).Offset(offset).Limit(f.Limit).Find(&items).Error
	return items, total, err
}

func (r *ContributionRepo) FindByID(id string) (*model.Contribution, error) {
	var c model.Contribution
	if err := r.db.First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ContributionRepo) Create(c *model.Contribution) error {
	return r.db.Create(c).Error
}

func (r *ContributionRepo) Update(c *model.Contribution) error {
	return r.db.Save(c).Error
}

func (r *ContributionRepo) FindAll() ([]model.Contribution, error) {
	var items []model.Contribution
	err := r.db.Order(`"createdAt" desc`).Find(&items).Error
	return items, err
}

func (r *ContributionRepo) SumVerified() (int64, error) {
	var result struct{ Sum int64 }
	err := r.db.Model(&model.Contribution{}).
		Where("status = ?", model.StatusVerified).
		Select("COALESCE(SUM(amount), 0) as sum").Scan(&result).Error
	return result.Sum, err
}

func (r *ContributionRepo) FindChartData(since time.Time) ([]model.Contribution, error) {
	var items []model.Contribution
	err := r.db.Where(`status = ? AND "verifiedAt" >= ?`, model.StatusVerified, since).
		Select(`amount, "verifiedAt"`).Find(&items).Error
	return items, err
}

func (r *ContributionRepo) Leaderboard() ([]model.LeaderboardItem, error) {
	var items []model.LeaderboardItem
	err := r.db.Model(&model.Contribution{}).
		Where("status = ?", model.StatusVerified).
		Select("name, COALESCE(SUM(amount), 0) as total").
		Group("name").
		Order("total DESC").
		Scan(&items).Error
	return items, err
}
