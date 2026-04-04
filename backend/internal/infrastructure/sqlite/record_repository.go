package sqlite

import (
	"github.com/Abhaykauts/LedgeGuard/backend/internal/domain"
	"gorm.io/gorm"
)

type recordRepository struct {
	db *gorm.DB
}

// NewRecordRepository creates a new SQLite record repository
func NewRecordRepository(db *gorm.DB) domain.RecordRepository {
	return &recordRepository{db: db}
}

func (r *recordRepository) Create(record *domain.Record) error {
	return r.db.Create(record).Error
}

func (r *recordRepository) GetByID(id uint) (*domain.Record, error) {
	var record domain.Record
	if err := r.db.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *recordRepository) Update(record *domain.Record) error {
	return r.db.Save(record).Error
}

func (r *recordRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Record{}, id).Error
}

func (r *recordRepository) List(filter domain.RecordFilter) ([]domain.Record, error) {
	var records []domain.Record
	query := r.db.Model(&domain.Record{})

	if filter.StartDate != nil {
		query = query.Where("date >= ?", filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("date <= ?", filter.EndDate)
	}
	if filter.Type != nil {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Category != nil {
		query = query.Where("category = ?", filter.Category)
	}

	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
