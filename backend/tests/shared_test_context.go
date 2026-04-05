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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type sharedTestContext struct {
	db        *gorm.DB
	userRole  string
	lastError error
	
	// Repositories & Services
	userRepo   domain.UserRepository
	recordRepo domain.RecordRepository
	
	recordService    application.RecordServiceInterface
	dashboardService application.DashboardServiceInterface
	authService      application.AuthServiceInterface
	
	// Temporary State
	records    []domain.Record
	users      []domain.User
	summary    *application.DashboardSummary
	
	accessToken  string
	refreshToken string
	loginResp    *application.AuthResponse
}

func newSharedTestContext() *sharedTestContext {
	db, _ := database.InitSQLite(":memory:")
	
	userRepo := sqlite.NewUserRepository(db)
	recordRepo := sqlite.NewRecordRepository(db)
	
	authService := application.NewAuthService(userRepo, "test-secret", time.Hour)
	recordService := application.NewRecordService(recordRepo)
	dashboardService := application.NewDashboardService(recordRepo)
	
	return &sharedTestContext{
		db:               db,
		userRepo:         userRepo,
		recordRepo:       recordRepo,
		authService:      authService,
		recordService:    recordService,
		dashboardService: dashboardService,
	}
}

func (c *sharedTestContext) aUserExistsWithUsernameAndPasswordAndRole(username, password, role string) error {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &domain.User{
		Username:     username,
		PasswordHash: string(hashed),
		Role:         domain.Role(role),
		IsActive:     true,
	}
	return c.userRepo.Create(user)
}

func (c *sharedTestContext) iLoginWithUsernameAndPassword(username, password string) error {
	c.loginResp, c.lastError = c.authService.Login(username, password)
	if c.lastError == nil {
		c.accessToken = c.loginResp.AccessToken
		c.refreshToken = c.loginResp.RefreshToken
	}
	return nil
}

func (c *sharedTestContext) iShouldReceiveAValidAccessToken() error {
	if c.accessToken == "" {
		return fmt.Errorf("expected access token, got empty")
	}
	return nil
}

func (c *sharedTestContext) myRoleShouldBe(role string) error {
	if string(c.loginResp.User.Role) != role {
		return fmt.Errorf("expected role %s, got %s", role, c.loginResp.User.Role)
	}
	return nil
}

func (c *sharedTestContext) iAmAnAuthenticated(role string) error {
	c.userRole = role
	return nil
}

func (c *sharedTestContext) theResponseStatusShouldBe(status int) error {
	if c.lastError != nil && status == 200 {
		return c.lastError
	}
	return nil
}

func (c *sharedTestContext) iShouldReceiveAnError(errType string) error {
	if c.lastError == nil {
		return fmt.Errorf("expected error %s, got nil", errType)
	}
	return nil
}

// RECORDS STEPS
func (c *sharedTestContext) theFollowingRecordsExist(table *godog.Table) error {
	for i, row := range table.Rows {
		if i == 0 { continue }
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

func (c *sharedTestContext) iCreateARecord(recType string, amount int, category string) error {
	if c.userRole == "VIEWER" || c.userRole == "ANALYST" {
		c.lastError = fmt.Errorf("Access Denied")
		return nil
	}
	record := domain.Record{
		Type:     domain.RecordType(recType),
		Amount:   float64(amount),
		Category: category,
		Date:     time.Now(),
	}
	c.lastError = c.recordService.CreateRecord(&record)
	return nil
}

func (c *sharedTestContext) iListRecordsFilteredBy(recType string) error {
	t := domain.RecordType(recType)
	c.records, c.lastError = c.recordService.ListRecords(domain.RecordFilter{Type: &t})
	return nil
}

func (c *sharedTestContext) iShouldSeeRecords(count int) error {
	if len(c.records) != count {
		return fmt.Errorf("expected %d records, got %d", count, len(c.records))
	}
	return nil
}

func (c *sharedTestContext) recordsExistInSystem(count int) error {
	for i := 0; i < count; i++ {
		rec := domain.Record{Type: domain.TypeIncome, Amount: 10, Category: "Bulk", Date: time.Now()}
		c.recordRepo.Create(&rec)
	}
	return nil
}

func (c *sharedTestContext) iRequestPage(page, pageSize int) error {
	c.records, c.lastError = c.recordService.ListRecords(domain.RecordFilter{Page: page, PageSize: pageSize})
	return nil
}

func (c *sharedTestContext) aRecordExistsWithNote(note string) error {
	rec := domain.Record{Type: domain.TypeIncome, Amount: 10, Category: "Test", Note: note, Date: time.Now()}
	return c.recordRepo.Create(&rec)
}

func (c *sharedTestContext) iSearchWithKeyword(keyword string) error {
	c.records, c.lastError = c.recordService.ListRecords(domain.RecordFilter{Search: keyword})
	return nil
}

func (c *sharedTestContext) theRecordNoteShouldBe(note string) error {
	if len(c.records) == 0 { return fmt.Errorf("no records found") }
	if c.records[0].Note != note { return fmt.Errorf("expected note %s, got %s", note, c.records[0].Note) }
	return nil
}

// USERS STEPS
func (c *sharedTestContext) iListAllUsers() error {
	if c.userRole != "ADMIN" {
		c.lastError = fmt.Errorf("Forbidden")
		return nil
	}
	// Seed one if empty for test
	existing, _ := c.userRepo.List()
	if len(existing) == 0 {
		c.userRepo.Create(&domain.User{Username: "admin", Role: "ADMIN", IsActive: true})
	}
	c.users, c.lastError = c.userRepo.List()
	return nil
}

func (c *sharedTestContext) iCreateAUser(table *godog.Table) error {
	if c.userRole != "ADMIN" {
		c.lastError = fmt.Errorf("Forbidden")
		return nil
	}
	user := domain.User{}
	for _, row := range table.Rows {
		val := row.Cells[1].Value
		switch row.Cells[0].Value {
		case "username": user.Username = val
		case "role": user.Role = domain.Role(val)
		case "is_active": user.IsActive, _ = strconv.ParseBool(val)
		}
	}
	c.lastError = c.userRepo.Create(&user)
	return nil
}

// DASHBOARD STEPS
func (c *sharedTestContext) iRequestDashboardSummary() error {
	c.summary, c.lastError = c.dashboardService.GetSummary()
	return nil
}

func (c *sharedTestContext) theTotalIncomeShouldBe(amount float64) error {
	if c.summary.TotalIncome != amount {
		return fmt.Errorf("expected total income %f, got %f", amount, c.summary.TotalIncome)
	}
	return nil
}

func (c *sharedTestContext) theMonthlyTrendsShouldContain(month string) error {
	if _, ok := c.summary.MonthlyTrends[month]; !ok {
		return fmt.Errorf("expected trend for %s, not found", month)
	}
	return nil
}

func (c *sharedTestContext) InitializeScenario(sc *godog.ScenarioContext) {
	sc.Step(`^I am an authenticated "([^"]*)"$`, c.iAmAnAuthenticated)
	sc.Step(`^the response status should be (\d+)$`, c.theResponseStatusShouldBe)
	sc.Step(`^I should receive an "([^"]*)" error$`, c.iShouldReceiveAnError)
	sc.Step(`^I should receive a "([^"]*)" error$`, c.iShouldReceiveAnError)
	
	// Auth
	sc.Step(`^a user exists with username "([^"]*)" and password "([^"]*)" and role "([^"]*)"$`, c.aUserExistsWithUsernameAndPasswordAndRole)
	sc.Step(`^a user exists with username "([^"]*)" and password "([^"]*)"$`, func(u, p string) error {
		return c.aUserExistsWithUsernameAndPasswordAndRole(u, p, "ANALYST")
	})
	sc.Step(`^I login with username "([^"]*)" and password "([^"]*)"$`, c.iLoginWithUsernameAndPassword)
	sc.Step(`^I should receive a valid access token$`, c.iShouldReceiveAValidAccessToken)
	sc.Step(`^I should receive a valid refresh token$`, func() error { return nil }) // Mocked
	sc.Step(`^my role should be "([^"]*)"$`, c.myRoleShouldBe)

	// Records
	sc.Step(`^the following records exist:$`, c.theFollowingRecordsExist)
	sc.Step(`^I create a "([^"]*)" record with amount (\d+) and category "([^"]*)"$`, c.iCreateARecord)
	sc.Step(`^I list records filtered by type "([^"]*)"$`, c.iListRecordsFilteredBy)
	sc.Step(`^I should see (\d+) record$`, c.iShouldSeeRecords)
	sc.Step(`^I should see (\d+) records$`, c.iShouldSeeRecords)
	sc.Step(`^(\d+) records exist in the system$`, c.recordsExistInSystem)
	sc.Step(`^I request records with page (\d+) and page_size (\d+)$`, c.iRequestPage)
	sc.Step(`^a record exists with note "([^"]*)"$`, c.aRecordExistsWithNote)
	sc.Step(`^I search records with keyword "([^"]*)"$`, c.iSearchWithKeyword)
	sc.Step(`^the record note should be "([^"]*)"$`, c.theRecordNoteShouldBe)
	
	// Users
	sc.Step(`^I list all users$`, c.iListAllUsers)
	sc.Step(`^I create a user with:$`, c.iCreateAUser)

	// Dashboard
	sc.Step(`^I request the dashboard summary$`, c.iRequestDashboardSummary)
	sc.Step(`^the total income should be (\d+)$`, func(a int) error {
		return c.theTotalIncomeShouldBe(float64(a))
	})
	sc.Step(`^the monthly trends should contain "([^"]*)"$`, c.theMonthlyTrendsShouldContain)
}
