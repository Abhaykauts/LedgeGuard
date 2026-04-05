package tests

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name:                "LedgeGuard Functional Tests",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"../features"},
			Output: colors.Colored(os.Stdout),
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	sharedCtx := newSharedTestContext()
	sharedCtx.InitializeScenario(sc)
}
