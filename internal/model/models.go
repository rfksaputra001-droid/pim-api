// pim-api-go/internal/model/models.go
package model

import "time"

type ContributionStatus string

const (
	StatusPending  ContributionStatus = "PENDING"
	StatusVerified ContributionStatus = "VERIFIED"
	StatusRejected ContributionStatus = "REJECTED"
)

type ExpenseCategory string

const (
	CategoryOperasional ExpenseCategory = "OPERASIONAL"
	CategoryKegiatan    ExpenseCategory = "KEGIATAN"
	CategoryPeralatan   ExpenseCategory = "PERALATAN"
	CategoryLainLain    ExpenseCategory = "LAIN_LAIN"
)

type AdminRole string

const (
	RoleAdmin      AdminRole = "ADMIN"
	RoleSuperAdmin AdminRole = "SUPER_ADMIN"
)

type Contribution struct {
	ID              string             `json:"id"              gorm:"primaryKey;type:text"`
	Name            string             `json:"name"`
	Phone           string             `json:"phone"`
	Amount          int64              `json:"amount"`
	Notes           *string            `json:"notes"`
	Status          ContributionStatus `json:"status"          gorm:"default:'PENDING'"`
	ProofImageURL   string             `json:"proofImageUrl"   gorm:"column:proof_image_url"`
	RejectionReason *string            `json:"rejectionReason"`
	CreatedAt       time.Time          `json:"createdAt"`
	VerifiedAt      *time.Time         `json:"verifiedAt"`
	VerifiedBy      *string            `json:"verifiedBy"`
}

type Expense struct {
	ID              string          `json:"id"              gorm:"primaryKey;type:text"`
	Date            time.Time       `json:"date"`
	Description     string          `json:"description"`
	Category        ExpenseCategory `json:"category"`
	Amount          int64           `json:"amount"`
	ReceiptImageURL string          `json:"receiptImageUrl" gorm:"column:receipt_image_url"`
	CreatedBy       string          `json:"createdBy"`
	CreatedAt       time.Time       `json:"createdAt"`
}

type Admin struct {
	ID           string     `json:"id"        gorm:"primaryKey;type:text"`
	Username     string     `json:"username"  gorm:"uniqueIndex"`
	PasswordHash string     `json:"-"`
	Role         AdminRole  `json:"role"      gorm:"default:'ADMIN'"`
	CreatedAt    time.Time  `json:"createdAt"`
	LastLogin    *time.Time `json:"lastLogin"`
}
