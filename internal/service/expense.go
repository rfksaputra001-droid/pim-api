package service

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"pim-api-go/internal/model"
	"pim-api-go/internal/repository"
)

type ExpenseService struct{ repo *repository.ExpenseRepo }

func NewExpenseService(repo *repository.ExpenseRepo) *ExpenseService {
	return &ExpenseService{repo}
}

type CreateExpenseInput struct {
	Tanggal    time.Time
	Keterangan string
	Kategori   model.ExpenseCategory
	Nominal    int64
	FileData   []byte
	MimeType   string
	CreatedBy  string
}

var validCategories = map[model.ExpenseCategory]bool{
	model.CategoryOperasional: true,
	model.CategoryKegiatan:    true,
	model.CategoryPeralatan:   true,
	model.CategoryLainLain:    true,
}

func (s *ExpenseService) Create(input CreateExpenseInput, uploadFn func([]byte, string) (string, error)) (*model.Expense, error) {
	if input.FileData == nil {
		return nil, fmt.Errorf("Foto nota wajib diupload")
	}
	if input.Nominal < 1 {
		return nil, fmt.Errorf("Nominal tidak valid")
	}
	if !validCategories[input.Kategori] {
		return nil, fmt.Errorf("Kategori tidak valid")
	}
	if input.Keterangan == "" {
		return nil, fmt.Errorf("Keterangan wajib diisi")
	}
	if input.Tanggal.IsZero() {
		return nil, fmt.Errorf("Tanggal wajib diisi")
	}

	receiptURL, err := uploadFn(input.FileData, "kas-rt/nota-pengeluaran")
	if err != nil {
		return nil, fmt.Errorf("Gagal mengupload foto nota")
	}

	expense := &model.Expense{
		ID:              uuid.NewString(),
		Date:            input.Tanggal,
		Description:     input.Keterangan,
		Category:        input.Kategori,
		Amount:          input.Nominal,
		ReceiptImageURL: receiptURL,
		CreatedBy:       input.CreatedBy,
		CreatedAt:       time.Now(),
	}
	if err := s.repo.Create(expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) ListAdmin(page int, category string) (*ListResult, error) {
	const limit = 20
	if page < 1 {
		page = 1
	}
	items, total, err := s.repo.FindAdmin(page, limit, category)
	if err != nil {
		return nil, err
	}
	return &ListResult{
		Data:       items,
		Total:      total,
		Page:       page,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}, nil
}

func (s *ExpenseService) ListPublic(page int, category string) (*ListResult, error) {
	const limit = 20
	if page < 1 {
		page = 1
	}
	items, total, err := s.repo.FindPublic(page, limit, category)
	if err != nil {
		return nil, err
	}
	return &ListResult{
		Data:       items,
		Total:      total,
		Page:       page,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}, nil
}

func (s *ExpenseService) GetReceipt(id string) (string, error) {
	e, err := s.repo.FindByID(id)
	if err != nil {
		return "", fmt.Errorf("Pengeluaran tidak ditemukan")
	}
	return e.ReceiptImageURL, nil
}
