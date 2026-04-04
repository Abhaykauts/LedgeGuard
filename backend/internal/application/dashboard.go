package application

// DashboardSummary represents the aggregated financial data
type DashboardSummary struct {
	TotalIncome    float64            `json:"total_income"`
	TotalExpenses  float64            `json:"total_expenses"`
	NetBalance     float64            `json:"net_balance"`
	CategoryTotals map[string]float64 `json:"category_totals"`
}

// DashboardServiceInterface defines the analytics use cases
type DashboardServiceInterface interface {
	GetSummary() (*DashboardSummary, error)
}
