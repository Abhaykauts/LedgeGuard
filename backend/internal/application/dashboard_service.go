package application

import (
	"github.com/Abhaykauts/LedgeGuard/backend/internal/domain"
)

type dashboardService struct {
	recordRepo domain.RecordRepository
}

// NewDashboardService creates a new instance of DashboardService
func NewDashboardService(repo domain.RecordRepository) DashboardServiceInterface {
	return &dashboardService{recordRepo: repo}
}

func (s *dashboardService) GetSummary() (*DashboardSummary, error) {
	records, err := s.recordRepo.List(domain.RecordFilter{})
	if err != nil {
		return nil, err
	}

	summary := &DashboardSummary{
		CategoryTotals: make(map[string]float64),
	}

	for _, rec := range records {
		if rec.Type == domain.TypeIncome {
			summary.TotalIncome += rec.Amount
		} else if rec.Type == domain.TypeExpense {
			summary.TotalExpenses += rec.Amount
		}

		summary.CategoryTotals[rec.Category] += rec.Amount
	}

	summary.NetBalance = summary.TotalIncome - summary.TotalExpenses

	return summary, nil
}
