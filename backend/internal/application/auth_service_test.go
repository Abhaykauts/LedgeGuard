package application_test

import (
	"time"

	"github.com/Abhaykauts/LedgeGuard/backend/internal/application"
	"github.com/Abhaykauts/LedgeGuard/backend/internal/domain"
	"github.com/Abhaykauts/LedgeGuard/backend/internal/infrastructure/sqlite"
	"github.com/Abhaykauts/LedgeGuard/backend/pkg/database"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var _ = Describe("AuthService", func() {
	var (
		db      *gorm.DB
		repo    domain.UserRepository
		service application.AuthServiceInterface
		secret  = "test-secret"
	)

	BeforeEach(func() {
		db, _ = database.InitSQLite(":memory:")
		repo = sqlite.NewUserRepository(db)
		service = application.NewAuthService(repo, secret, time.Hour)
	})

	Context("Login Workflow", func() {
		BeforeEach(func() {
			hashed, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
			repo.Create(&domain.User{
				Username:     "tester",
				PasswordHash: string(hashed),
				Role:         domain.RoleAdmin,
				IsActive:     true,
			})
		})

		It("should successfully login with valid credentials", func() {
			resp, err := service.Login("tester", "pass123")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.AccessToken).NotTo(BeEmpty())
			Expect(resp.User.Username).To(Equal("tester"))
		})

		It("should fail with invalid password", func() {
			_, err := service.Login("tester", "wrong-pass")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("invalid credentials"))
		})

		It("should fail for non-existent user", func() {
			_, err := service.Login("non-existent", "pass123")
			Expect(err).To(HaveOccurred())
		})
	})
})
