package application

import "github.com/Abhaykauts/LedgeGuard/backend/internal/domain"

// RecordServiceInterface defines the financial record use cases
type RecordServiceInterface interface {
	CreateRecord(record *domain.Record) error
	GetRecord(id uint) (*domain.Record, error)
	UpdateRecord(record *domain.Record) error
	DeleteRecord(id uint) error
	ListRecords(filter domain.RecordFilter) ([]domain.Record, error)
}
