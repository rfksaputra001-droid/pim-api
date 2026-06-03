package service

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"pim-api-go/internal/model"
	"pim-api-go/internal/repository"
	"pim-api-go/pkg/crypto"
	"pim-api-go/pkg/walink"
)

type ContributionService struct {
	repo        *repository.ContributionRepo
	expenseRepo *repository.ExpenseRepo
}

func NewContributionService(repo *repository.ContributionRepo, expenseRepo *repository.ExpenseRepo) *ContributionService {
	return &ContributionService{repo, expenseRepo}
}

type DashboardResult struct {
	TotalIncome   int64       `json:"totalIncome"`
	TotalExpenses int64       `json:"totalExpenses"`
	Balance       int64       `json:"balance"`
	ChartData     []ChartItem `json:"chartData"`
}

type ChartItem struct {
	Bulan       string `json:"bulan"`
	Pemasukan   int64  `json:"pemasukan"`
	Pengeluaran int64  `json:"pengeluaran"`
}

var monthNames = []string{"Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"}

func (s *ContributionService) Dashboard() (*DashboardResult, error) {
	totalIncome, err := s.repo.SumVerified()
	if err != nil {
		return nil, err
	}
	totalExpenses, err := s.expenseRepo.SumAll()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	sixMonthsAgo := time.Date(now.Year(), now.Month()-5, 1, 0, 0, 0, 0, time.UTC)

	contributions, _ := s.repo.FindChartData(sixMonthsAgo)
	expenses, _ := s.expenseRepo.FindChartData(sixMonthsAgo)

	chartData := make([]ChartItem, 6)
	for i := 5; i >= 0; i-- {
		d := time.Date(now.Year(), now.Month()-time.Month(i), 1, 0, 0, 0, 0, time.UTC)
		m, y := int(d.Month()), d.Year()
		idx := 5 - i

		var income, expense int64
		for _, c := range contributions {
			if c.VerifiedAt != nil && int(c.VerifiedAt.Month()) == m && c.VerifiedAt.Year() == y {
				income += c.Amount
			}
		}
		for _, e := range expenses {
			if int(e.Date.Month()) == m && e.Date.Year() == y {
				expense += e.Amount
			}
		}

		chartData[idx] = ChartItem{
			Bulan:       fmt.Sprintf("%s %d", monthNames[m-1], y),
			Pemasukan:   income,
			Pengeluaran: expense,
		}
	}

	return &DashboardResult{
		TotalIncome:   totalIncome,
		TotalExpenses: totalExpenses,
		Balance:       totalIncome - totalExpenses,
		ChartData:     chartData,
	}, nil
}

type ListResult struct {
	Data       any   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	TotalPages int   `json:"totalPages"`
}

func (s *ContributionService) ListPublic(year, month, page int) (*ListResult, error) {
	const limit = 10
	if page < 1 {
		page = 1
	}
	items, total, err := s.repo.FindPublic(year, month, page, limit)
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

type SubmitInput struct {
	Nama        string
	NoTelepon   string
	Nominal     int64
	Catatan     string
	IsAnonymous bool
	FileData    []byte
	MimeType    string
}

func (s *ContributionService) Submit(input SubmitInput, uploadFn func([]byte, string) (string, error)) error {
	if input.FileData == nil {
		return fmt.Errorf("Bukti transfer wajib diupload")
	}
	if input.Nominal < 1000 {
		return fmt.Errorf("Nominal minimal Rp 1.000")
	}
	if len(input.NoTelepon) < 9 {
		return fmt.Errorf("Nomor telepon tidak valid")
	}

	name := input.Nama
	if input.IsAnonymous {
		name = "Hamba Allah"
	}
	if name == "" {
		return fmt.Errorf("Nama wajib diisi")
	}

	proofURL, err := uploadFn(input.FileData, "kas-rt/bukti-iuran")
	if err != nil {
		return fmt.Errorf("Gagal mengupload bukti transfer")
	}

	encPhone, err := crypto.Encrypt(input.NoTelepon)
	if err != nil {
		return err
	}

	var notes *string
	if input.Catatan != "" {
		notes = &input.Catatan
	}

	return s.repo.Create(&model.Contribution{
		ID:            uuid.NewString(),
		Name:          name,
		Phone:         encPhone,
		Amount:        input.Nominal,
		Notes:         notes,
		ProofImageURL: proofURL,
	})
}

func (s *ContributionService) ListAdmin(status, sortBy, sortDir string, page int) (*ListResult, error) {
	const limit = 20
	if page < 1 {
		page = 1
	}
	f := repository.ContributionFilter{
		Status:  model.ContributionStatus(status),
		SortBy:  sortBy,
		SortDir: sortDir,
		Page:    page,
		Limit:   limit,
	}
	items, total, err := s.repo.FindAdmin(f)
	if err != nil {
		return nil, err
	}
	for i := range items {
		if decrypted, err := crypto.Decrypt(items[i].Phone); err == nil {
			items[i].Phone = decrypted
		} else {
			items[i].Phone = "(enkripsi error)"
		}
	}
	return &ListResult{
		Data:       items,
		Total:      total,
		Page:       page,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}, nil
}

func (s *ContributionService) Verify(id, verifiedBy string) (string, error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return "", fmt.Errorf("Iuran tidak ditemukan")
	}
	if c.Status != model.StatusPending {
		return "", fmt.Errorf("Iuran sudah diproses sebelumnya")
	}

	now := time.Now()
	c.Status = model.StatusVerified
	c.VerifiedAt = &now
	c.VerifiedBy = &verifiedBy
	if err := s.repo.Update(c); err != nil {
		return "", err
	}

	income, _ := s.repo.SumVerified()
	expenses, _ := s.expenseRepo.SumAll()
	balance := income - expenses

	phone, _ := crypto.Decrypt(c.Phone)
	return walink.Verified(phone, c.Name, c.Amount, balance), nil
}

func (s *ContributionService) Reject(id, reason string) (string, error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return "", fmt.Errorf("Iuran tidak ditemukan")
	}
	if c.Status != model.StatusPending {
		return "", fmt.Errorf("Iuran sudah diproses sebelumnya")
	}

	c.Status = model.StatusRejected
	c.RejectionReason = &reason
	if err := s.repo.Update(c); err != nil {
		return "", err
	}

	phone, _ := crypto.Decrypt(c.Phone)
	return walink.Rejected(phone, c.Name, c.Amount, reason), nil
}

func (s *ContributionService) ExportCSV() ([][]string, error) {
	items, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	fmtDate := func(t *time.Time) string {
		if t == nil {
			return ""
		}
		return t.Format("02/01/2006")
	}
	strOrEmpty := func(p *string) string {
		if p == nil {
			return ""
		}
		return *p
	}

	rows := [][]string{
		{"Tanggal Submit", "Nama", "Nominal (Rp)", "Status", "Catatan", "Diverifikasi Oleh", "Tanggal Verifikasi", "Alasan Tolak"},
	}
	for _, c := range items {
		rows = append(rows, []string{
			fmtDate(&c.CreatedAt),
			c.Name,
			fmt.Sprintf("%d", c.Amount),
			string(c.Status),
			strOrEmpty(c.Notes),
			strOrEmpty(c.VerifiedBy),
			fmtDate(c.VerifiedAt),
			strOrEmpty(c.RejectionReason),
		})
	}
	return rows, nil
}
