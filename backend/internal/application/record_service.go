package application

import (
	"errors"

	"github.com/Abhaykauts/LedgeGuard/backend/internal/domain"
)

type recordService struct {
	repo domain.RecordRepository
}

// NewRecordService creates a new instance of RecordService
func NewRecordService(repo domain.RecordRepository) RecordServiceInterface {
	return &recordService{repo: repo}
}

func (s *recordService) CreateRecord(record *domain.Record) error {
	if record.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	return s.repo.Create(record)
}

func (s *recordService) GetRecord(id uint) (*domain.Record, error) {
	return s.repo.GetByID(id)
}

func (s *recordService) UpdateRecord(record *domain.Record) error {
	if record.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	return s.repo.Update(record)
}

func (s *recordService) DeleteRecord(id uint) error {
	return s.repo.Delete(id)
}

func (s *recordService) ListRecords(filter domain.RecordFilter) ([]domain.Record, error) {
	return s.repo.List(filter)
}
