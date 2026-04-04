package tests

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

func TestMain(m *testing.M) {
	status := godog.TestSuite{
		Name:                "LedgeGuard Functional Tests",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"../features"},
			Output: colors.Colored(os.Stdout),
		},
	}.Run()

	if status != 0 {
		os.Exit(status)
	}

	os.Exit(m.Run())
}

func InitializeScenario(sc *godog.ScenarioContext) {
	authCtx := newAuthTestContext()
	authCtx.InitializeScenario(sc)

	recordCtx := newRecordTestContext()
	recordCtx.InitializeScenario(sc)

	dashboardCtx := newDashboardTestContext()
	dashboardCtx.InitializeScenario(sc)
}
