package domain

import (
	"time"

	"gorm.io/gorm"
)

// RecordType defines if it's income or expense
type RecordType string

const (
	TypeIncome  RecordType = "INCOME"
	TypeExpense RecordType = "EXPENSE"
)

// Record represents a single financial entry
type Record struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Amount      float64        `json:"amount"`
	Type        RecordType     `json:"type"`
	Category    string         `json:"category"`
	Date        time.Time      `json:"date"`
	Note        string         `json:"note"`
	CreatedBy   uint           `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// RecordRepository defines the contract for record persistence
type RecordRepository interface {
	Create(record *Record) error
	GetByID(id uint) (*Record, error)
	Update(record *Record) error
	Delete(id uint) error
	List(filter RecordFilter) ([]Record, error)
}

// RecordFilter provides criteria for listing records
type RecordFilter struct {
	StartDate *time.Time  `form:"start_date" time_format:"2006-01-02"`
	EndDate   *time.Time  `form:"end_date" time_format:"2006-01-02"`
	Type      *RecordType `form:"type"`
	Category  *string     `form:"category"`
	Page      int         `form:"page,default=1"`
	PageSize  int         `form:"page_size,default=10"`
	Search    string      `form:"search"`
}
