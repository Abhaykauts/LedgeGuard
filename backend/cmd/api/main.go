package main

import (
	"log"

	"github.com/Abhaykauts/LedgeGuard/backend/internal/application"
	"github.com/Abhaykauts/LedgeGuard/backend/internal/domain"
	infra_config "github.com/Abhaykauts/LedgeGuard/backend/internal/infrastructure/config"
	"github.com/Abhaykauts/LedgeGuard/backend/internal/infrastructure/sqlite"
	"github.com/Abhaykauts/LedgeGuard/backend/internal/interface/http"
	"github.com/Abhaykauts/LedgeGuard/backend/pkg/auth"
	"github.com/Abhaykauts/LedgeGuard/backend/pkg/database"
)

func main() {
	// 1. Load Config
	cfg := infra_config.LoadConfig()

	// 2. Init DB
	db, err := database.InitSQLite(cfg.DBPath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 3. Setup Layers
	userRepo := sqlite.NewUserRepository(db)
	recordRepo := sqlite.NewRecordRepository(db)

	authService := application.NewAuthService(userRepo, cfg.JWTSecret, cfg.TokenDuration)
	recordService := application.NewRecordService(recordRepo)

	authHandler := http.NewAuthHandler(authService)
	recordHandler := http.NewRecordHandler(recordService)

	// 4. Seed initial Admin (for testing)
	seedAdmin(userRepo)

	// 5. Setup Router
	routerCfg := http.RouterConfig{
		AuthHandler:   authHandler,
		RecordHandler: recordHandler,
		JWTSecret:     cfg.JWTSecret,
	}
	r := http.NewRouter(routerCfg)

	// 6. Start Server
	log.Printf("LedgeGuard API starting on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func seedAdmin(repo domain.UserRepository) {
	username := "admin"
	if _, err := repo.GetByUsername(username); err != nil {
		hash, _ := auth.HashPassword("admin123")
		admin := &domain.User{
			Username:     username,
			PasswordHash: hash,
			Role:         domain.RoleAdmin,
			IsActive:     true,
		}
		repo.Create(admin)
		log.Println("Seeded initial admin user: admin / admin123")
	}
}
