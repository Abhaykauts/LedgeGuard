package application_test

import (
	"time"

	"github.com/Abhaykauts/LedgeGuard/backend/internal/application"
	"github.com/Abhaykauts/LedgeGuard/backend/internal/domain"
	"github.com/Abhaykauts/LedgeGuard/backend/internal/infrastructure/sqlite"
	"github.com/Abhaykauts/LedgeGuard/backend/pkg/database"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("DashboardService", func() {
	var (
		db      *gorm.DB
		repo    domain.RecordRepository
		service application.DashboardServiceInterface
	)

	BeforeEach(func() {
		db, _ = database.InitSQLite(":memory:")
		repo = sqlite.NewRecordRepository(db)
		service = application.NewDashboardService(repo)
	})

	Context("Summary Generation", func() {
		It("should correctly aggregate totals and trends", func() {
			// Seed data across multiple days/weeks
			t1, _ := time.Parse("2006-01-02", "2026-04-01") // Week 14
			t2, _ := time.Parse("2006-01-02", "2026-04-08") // Week 15
			
			repo.Create(&domain.Record{Type: domain.TypeIncome, Amount: 1000, Note: "W14 Pay", Date: t1})
			repo.Create(&domain.Record{Type: domain.TypeExpense, Amount: 200, Note: "W14 Food", Date: t1})
			repo.Create(&domain.Record{Type: domain.TypeIncome, Amount: 500, Note: "W15 Gift", Date: t2})

			summary, err := service.GetSummary()
			Expect(err).NotTo(HaveOccurred())
			Expect(summary.TotalIncome).To(Equal(float64(1500)))
			Expect(summary.TotalExpenses).To(Equal(float64(200)))
			Expect(summary.NetBalance).To(Equal(float64(1300)))

			// Weekly Trends (ISO-8601 keys)
			Expect(summary.WeeklyTrends).To(HaveKey("2026-W14"))
			Expect(summary.WeeklyTrends["2026-W14"]).To(Equal(float64(800))) // 1000 - 200
			Expect(summary.WeeklyTrends).To(HaveKey("2026-W15"))
			Expect(summary.WeeklyTrends["2026-W15"]).To(Equal(float64(500)))
		})
	})
})
