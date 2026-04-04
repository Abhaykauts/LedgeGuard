package auth_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Abhaykauts/LedgeGuard/backend/pkg/auth"
)

var _ = Describe("Auth Utils", func() {
	Context("Password Hashing", func() {
		It("should hash and verify passwords correctly", func() {
			password := "securePassword123"
			hash, err := auth.HashPassword(password)
			Expect(err).NotTo(HaveOccurred())
			Expect(hash).NotTo(Equal(password))

			isValid := auth.CheckPasswordHash(password, hash)
			Expect(isValid).To(BeTrue())

			isInvalid := auth.CheckPasswordHash("wrongPassword", hash)
			Expect(isInvalid).To(BeFalse())
		})
	})

	Context("JWT Tokens", func() {
		It("should generate and validate tokens correctly", func() {
			secret := "test-secret"
			userID := uint(1)
			role := "ADMIN"

			token, err := auth.GenerateToken(userID, role, secret, time.Hour)
			Expect(err).NotTo(HaveOccurred())
			Expect(token).NotTo(BeEmpty())

			claims, err := auth.ValidateToken(token, secret)
			Expect(err).NotTo(HaveOccurred())
			Expect(claims.UserID).To(Equal(userID))
			Expect(claims.Role).To(Equal(role))
		})

		It("should fail for invalid secrets", func() {
			token, _ := auth.GenerateToken(1, "ADMIN", "secret1", time.Hour)
			_, err := auth.ValidateToken(token, "wrong-secret")
			Expect(err).To(HaveOccurred())
		})

		It("should fail for expired tokens", func() {
			token, _ := auth.GenerateToken(1, "ADMIN", "secret", -time.Hour)
			_, err := auth.ValidateToken(token, "secret")
			Expect(err).To(HaveOccurred())
		})
	})
})
