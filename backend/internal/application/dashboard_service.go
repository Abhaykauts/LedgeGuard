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
	// 1. Fetch ALL records for calculations (considering large datasets, this would be optimized)
	records, err := s.recordRepo.List(domain.RecordFilter{Page: 1, PageSize: 1000})
	if err != nil {
		return nil, err
	}

	summary := &DashboardSummary{
		CategoryTotals: make(map[string]float64),
		MonthlyTrends:  make(map[string]float64),
		RecentActivity: make([]domain.Record, 0),
	}

	for _, rec := range records {
		if rec.Type == domain.TypeIncome {
			summary.TotalIncome += rec.Amount
		} else if rec.Type == domain.TypeExpense {
			summary.TotalExpenses += rec.Amount
		}

		summary.CategoryTotals[rec.Category] += rec.Amount

		// Monthly Trend Calculation (Format: YYYY-MM)
		monthKey := rec.Date.Format("2006-01")
		if rec.Type == domain.TypeIncome {
			summary.MonthlyTrends[monthKey] += rec.Amount
		} else {
			summary.MonthlyTrends[monthKey] -= rec.Amount
		}
	}

	// 2. Fetch Recent Activity (explicitly use pagination for latest 5)
	recentFilter := domain.RecordFilter{Page: 1, PageSize: 5}
	summary.RecentActivity, _ = s.recordRepo.List(recentFilter)

	summary.NetBalance = summary.TotalIncome - summary.TotalExpenses

	return summary, nil
}
