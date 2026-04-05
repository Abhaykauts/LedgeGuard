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

var _ = Describe("RecordService", func() {
	var (
		db      *gorm.DB
		repo    domain.RecordRepository
		service application.RecordServiceInterface
	)

	BeforeEach(func() {
		db, _ = database.InitSQLite(":memory:")
		repo = sqlite.NewRecordRepository(db)
		service = application.NewRecordService(repo)
	})

	Describe("Creating Records", func() {
		It("should successfully create a valid record", func() {
			record := &domain.Record{
				Type:     domain.TypeIncome,
				Amount:   1000.50,
				Category: "Salary",
				Date:     time.Now(),
			}

			err := service.CreateRecord(record)
			Expect(err).NotTo(HaveOccurred())
			Expect(record.ID).NotTo(BeZero())
		})

		It("should fail if amount is zero or negative", func() {
			record := &domain.Record{
				Type:     domain.TypeExpense,
				Amount:   0,
				Category: "Coffee",
			}

			err := service.CreateRecord(record)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("greater than zero"))
		})
	})

	Describe("Listing & Filtering", func() {
		BeforeEach(func() {
			// Seed data
			repo.Create(&domain.Record{Type: domain.TypeIncome, Amount: 100, Note: "Project payment", Date: time.Now()})
			repo.Create(&domain.Record{Type: domain.TypeIncome, Amount: 200, Note: "Bonus", Date: time.Now()})
			repo.Create(&domain.Record{Type: domain.TypeExpense, Amount: 50, Note: "Lunch", Date: time.Now()})
		})

		It("should filter by keyword in notes", func() {
			filter := domain.RecordFilter{Search: "payment"}
			records, err := service.ListRecords(filter)
			Expect(err).NotTo(HaveOccurred())
			Expect(records).To(HaveLen(1))
			Expect(records[0].Note).To(Equal("Project payment"))
		})

		It("should handle pagination (limit and offset)", func() {
			filter := domain.RecordFilter{Page: 2, PageSize: 1}
			records, err := service.ListRecords(filter)
			Expect(err).NotTo(HaveOccurred())
			Expect(records).To(HaveLen(1))
			// Second record in order (usually IDs are predictable here)
			Expect(records[0].Amount).To(Equal(float64(200)))
		})
	})
})
