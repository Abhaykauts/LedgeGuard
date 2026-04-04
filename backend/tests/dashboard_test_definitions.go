package tests

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Abhaykauts/LedgeGuard/backend/internal/application"
	"github.com/Abhaykauts/LedgeGuard/backend/internal/domain"
	"github.com/Abhaykauts/LedgeGuard/backend/internal/infrastructure/sqlite"
	"github.com/Abhaykauts/LedgeGuard/backend/pkg/database"
	"github.com/cucumber/godog"
)

type dashboardTestContext struct {
	recordRepo       domain.RecordRepository
	dashboardService application.DashboardServiceInterface
	summary          *application.DashboardSummary
}

func newDashboardTestContext() *dashboardTestContext {
	db, _ := database.InitSQLite(":memory:")
	repo := sqlite.NewRecordRepository(db)
	service := application.NewDashboardService(repo)
	return &dashboardTestContext{
		recordRepo:       repo,
		dashboardService: service,
	}
}

func (c *dashboardTestContext) theFollowingFinancialRecordsExist(table *godog.Table) error {
	for i, row := range table.Rows {
		if i == 0 {
			continue
		} // skip header
		amount, _ := strconv.ParseFloat(row.Cells[1].Value, 64)
		record := domain.Record{
			Type:     domain.RecordType(row.Cells[0].Value),
			Amount:   amount,
			Category: row.Cells[2].Value,
			Date:     time.Now(),
		}
		c.recordRepo.Create(&record)
	}
	return nil
}

func (c *dashboardTestContext) iRequestTheDashboardSummary() error {
	var err error
	c.summary, err = c.dashboardService.GetSummary()
	return err
}

func (c *dashboardTestContext) theTotalIncomeShouldBe(expected float64) error {
	if c.summary.TotalIncome != expected {
		return fmt.Errorf("expected %f, got %f", expected, c.summary.TotalIncome)
	}
	return nil
}

func (c *dashboardTestContext) theTotalExpensesShouldBe(expected float64) error {
	if c.summary.TotalExpenses != expected {
		return fmt.Errorf("expected %f, got %f", expected, c.summary.TotalExpenses)
	}
	return nil
}

func (c *dashboardTestContext) theNetBalanceShouldBe(expected float64) error {
	if c.summary.NetBalance != expected {
		return fmt.Errorf("expected %f, got %f", expected, c.summary.NetBalance)
	}
	return nil
}

func (c *dashboardTestContext) theCategoryTotalShouldBe(category string, expected float64) error {
	total := c.summary.CategoryTotals[category]
	if total != expected {
		return fmt.Errorf("expected %f for category %s, got %f", expected, category, total)
	}
	return nil
}

func (c *dashboardTestContext) InitializeScenario(sc *godog.ScenarioContext) {
	sc.Step(`^the following financial records exist:$`, c.theFollowingFinancialRecordsExist)
	sc.Step(`^I request the dashboard summary$`, c.iRequestTheDashboardSummary)
	sc.Step(`^the total income should be (\d+.*)$`, func(a float64) error {
		return c.theTotalIncomeShouldBe(a)
	})
	sc.Step(`^the total expenses should be (\d+.*)$`, func(a float64) error {
		return c.theTotalExpensesShouldBe(a)
	})
	sc.Step(`^the net balance should be (\d+.*)$`, func(a float64) error {
		return c.theNetBalanceShouldBe(a)
	})
	sc.Step(`^the category "([^"]*)" total should be (\d+.*)$`, func(cat string, a float64) error {
		return c.theCategoryTotalShouldBe(cat, a)
	})
}
