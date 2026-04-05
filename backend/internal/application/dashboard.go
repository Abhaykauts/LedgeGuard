package application

import "github.com/Abhaykauts/LedgeGuard/backend/internal/domain"

// DashboardSummary represents the aggregated financial data
type DashboardSummary struct {
	TotalIncome    float64            `json:"total_income"`
	TotalExpenses  float64            `json:"total_expenses"`
	NetBalance     float64            `json:"net_balance"`
	CategoryTotals map[string]float64 `json:"category_totals"`
	RecentActivity []domain.Record    `json:"recent_activity"`
	MonthlyTrends  map[string]float64 `json:"monthly_trends"`
	WeeklyTrends   map[string]float64 `json:"weekly_trends"`
}

// DashboardServiceInterface defines the analytics use cases
type DashboardServiceInterface interface {
	GetSummary() (*DashboardSummary, error)
}
