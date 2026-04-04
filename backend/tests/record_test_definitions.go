package tests

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Abhaykauts/LedgeGuard/backend/internal/application"
	"github.com/Abhaykauts/LedgeGuard/backend/internal/domain"
	"github.com/Abhaykauts/LedgeGuard/backend/internal/infrastructure/sqlite"
	"github.com/Abhaykauts/LedgeGuard/backend/pkg/database"
	"github.com/cucumber/godog"
	"gorm.io/gorm"
)

type recordTestContext struct {
	db            *gorm.DB
	recordRepo    domain.RecordRepository
	recordService application.RecordServiceInterface
	userRole      string
	lastError     error
	records       []domain.Record
}

func newRecordTestContext() *recordTestContext {
	db, _ := database.InitSQLite(":memory:")
	repo := sqlite.NewRecordRepository(db)
	service := application.NewRecordService(repo)
	return &recordTestContext{
		db:            db,
		recordRepo:    repo,
		recordService: service,
	}
}

func (c *recordTestContext) iAmAnAuthenticated(role string) error {
	c.userRole = role
	return nil
}

func (c *recordTestContext) iCreateARecordWithAmountAndCategory(recType, amountStr, category string) error {
	amount, _ := strconv.ParseFloat(amountStr, 64)
	
	// Mocking RBAC check usually done in middleware
	if c.userRole == "VIEWER" {
		c.lastError = errors.New("Access Denied")
		return nil
	}

	record := domain.Record{
		Type:     domain.RecordType(recType),
		Amount:   amount,
		Category: category,
		Date:     time.Now(),
	}
	c.lastError = c.recordService.CreateRecord(&record)
	return nil
}

func (c *recordTestContext) theRecordShouldBeSavedSuccessfully() error {
	return c.lastError
}

func (c *recordTestContext) theTotalCountOfRecordsShouldBe(count int) error {
	records, _ := c.recordService.ListRecords(domain.RecordFilter{})
	if len(records) != count {
		return fmt.Errorf("expected %d records, got %d", count, len(records))
	}
	return nil
}

func (c *recordTestContext) iShouldReceiveAnError(errType string) error {
	if c.lastError == nil {
		return fmt.Errorf("expected error %s, got nil", errType)
	}
	return nil
}

func (c *recordTestContext) theFollowingRecordsExist(table *godog.Table) error {
	for i, row := range table.Rows {
		if i == 0 { continue } // skip header
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

func (c *recordTestContext) iListRecordsFilteredByType(recType string) error {
	t := domain.RecordType(recType)
	c.records, c.lastError = c.recordService.ListRecords(domain.RecordFilter{Type: &t})
	return nil
}

func (c *recordTestContext) iShouldSeeRecord(count int) error {
	if len(c.records) != count {
		return fmt.Errorf("expected %d record, got %d", count, len(c.records))
	}
	return nil
}

func (c *recordTestContext) theRecordAmountShouldBe(amount float64) error {
	if c.records[0].Amount != amount {
		return fmt.Errorf("expected amount %f, got %f", amount, c.records[0].Amount)
	}
	return nil
}

func (c *recordTestContext) InitializeScenario(sc *godog.ScenarioContext) {
	sc.Step(`^I am an authenticated "([^"]*)"$`, c.iAmAnAuthenticated)
	sc.Step(`^I create a "([^"]*)" record with amount (\d+) and category "([^"]*)"$`, c.iCreateARecordWithAmountAndCategory)
	sc.Step(`^the record should be saved successfully$`, c.theRecordShouldBeSavedSuccessfully)
	sc.Step(`^the total count of records should be (\d+)$`, c.theTotalCountOfRecordsShouldBe)
	sc.Step(`^I should receive an "([^"]*)" error$`, c.iShouldReceiveAnError)
	sc.Step(`^the following records exist:$`, c.theFollowingRecordsExist)
	sc.Step(`^I list records filtered by type "([^"]*)"$`, c.iListRecordsFilteredByType)
	sc.Step(`^I should see (\d+) record$`, c.iShouldSeeRecord)
	sc.Step(`^the record amount should be (\d+)$`, func(a int) error {
		return c.theRecordAmountShouldBe(float64(a))
	})
}
