package tests

import (
	"errors"
	"fmt"
	"time"

	"github.com/Abhaykauts/LedgeGuard/backend/internal/application"
	"github.com/Abhaykauts/LedgeGuard/backend/internal/domain"
	"github.com/Abhaykauts/LedgeGuard/backend/internal/infrastructure/sqlite"
	"github.com/Abhaykauts/LedgeGuard/backend/pkg/auth"
	"github.com/Abhaykauts/LedgeGuard/backend/pkg/database"
	"github.com/cucumber/godog"
	"gorm.io/gorm"
)

type authTestContext struct {
	db          *gorm.DB
	userRepo    domain.UserRepository
	authService application.AuthServiceInterface
	lastOutput  *application.AuthResponse
	lastError   error
	jwtSecret   string
}

func newAuthTestContext() *authTestContext {
	db, _ := database.InitSQLite(":memory:") // Use in-memory DB for tests
	repo := sqlite.NewUserRepository(db)
	secret := "test-secret"
	service := application.NewAuthService(repo, secret, time.Hour)
	return &authTestContext{
		db:          db,
		userRepo:    repo,
		authService: service,
		jwtSecret:   secret,
	}
}

func (c *authTestContext) aUserExistsWithUsernameAndPasswordAndRole(username, password, role string) error {
	hash, _ := auth.HashPassword(password)
	user := domain.User{
		Username:     username,
		PasswordHash: hash,
		Role:         domain.Role(role),
		IsActive:     true,
	}
	return c.userRepo.Create(&user)
}

func (c *authTestContext) iLoginWithUsernameAndPassword(username, password string) error {
	c.lastOutput, c.lastError = c.authService.Login(username, password)
	return nil
}

func (c *authTestContext) iShouldReceiveAValidAccessToken() error {
	if c.lastError != nil {
		return c.lastError
	}
	if c.lastOutput.AccessToken == "" {
		return errors.New("access token is empty")
	}
	_, err := auth.ValidateToken(c.lastOutput.AccessToken, c.jwtSecret)
	return err
}

func (c *authTestContext) iShouldReceiveAValidRefreshToken() error {
	if c.lastOutput.RefreshToken == "" {
		return errors.New("refresh token is empty")
	}
	return nil
}

func (c *authTestContext) myRoleShouldBe(role string) error {
	if string(c.lastOutput.User.Role) != role {
		return fmt.Errorf("expected role %s, got %s", role, c.lastOutput.User.Role)
	}
	return nil
}

func (c *authTestContext) iShouldReceiveAnUnauthorizedError() error {
	if c.lastError == nil {
		return errors.New("expected unathorized error, got nil")
	}
	if c.lastError.Error() != "invalid credentials" {
		return fmt.Errorf("expected 'invalid credentials' error, got '%v'", c.lastError)
	}
	return nil
}

func (c *authTestContext) iHaveAValidRefreshTokenForUser(username string) error {
	user, err := c.userRepo.GetByUsername(username)
	if err != nil {
		return err
	}
	token, _ := auth.GenerateToken(user.ID, string(user.Role), c.jwtSecret, time.Hour*24)
	c.lastOutput = &application.AuthResponse{RefreshToken: token}
	return nil
}

func (c *authTestContext) iRequestANewAccessTokenUsingTheRefreshToken() error {
	c.lastOutput, c.lastError = c.authService.RefreshToken(c.lastOutput.RefreshToken)
	return nil
}

func (c *authTestContext) InitializeScenario(sc *godog.ScenarioContext) {
	sc.Step(`^a user exists with username "([^"]*)" and password "([^"]*)" and role "([^"]*)"$`, c.aUserExistsWithUsernameAndPasswordAndRole)
	sc.Step(`^a user exists with username "([^"]*)" and password "([^"]*)"$`, func(u, p string) error {
		return c.aUserExistsWithUsernameAndPasswordAndRole(u, p, "VIEWER")
	})
	sc.Step(`^I login with username "([^"]*)" and password "([^"]*)"$`, c.iLoginWithUsernameAndPassword)
	sc.Step(`^I should receive a valid access token$`, c.iShouldReceiveAValidAccessToken)
	sc.Step(`^I should receive a valid refresh token$`, c.iShouldReceiveAValidRefreshToken)
	sc.Step(`^my role should be "([^"]*)"$`, c.myRoleShouldBe)
	sc.Step(`^I should receive an "Unauthorized" error$`, c.iShouldReceiveAnUnauthorizedError)
	sc.Step(`^I have a valid refresh token for user "([^"]*)"$`, c.iHaveAValidRefreshTokenForUser)
	sc.Step(`^I request a new access token using the refresh token$`, c.iRequestANewAccessTokenUsingTheRefreshToken)
	sc.Step(`^I should receive a new valid access token$`, c.iShouldReceiveAValidAccessToken)
	sc.Step(`^I should receive a new valid refresh token$`, c.iShouldReceiveAValidRefreshToken)
}

func (c *authTestContext) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		// Global setup
	})
}
