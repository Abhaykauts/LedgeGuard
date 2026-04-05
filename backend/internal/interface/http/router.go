package http

import (
	_ "github.com/Abhaykauts/LedgeGuard/backend/docs"
	"github.com/Abhaykauts/LedgeGuard/backend/internal/interface/http/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RouterConfig contains dependencies for the router
type RouterConfig struct {
	AuthHandler      *AuthHandler
	RecordHandler    *RecordHandler
	DashboardHandler *DashboardHandler
	JWTSecret        string
}

// NewRouter sets up the API routes
func NewRouter(cfg RouterConfig) *gin.Engine {
	r := gin.Default()

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Base API group
	api := r.Group("/api")
	{
		// Auth routes
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", cfg.AuthHandler.Login)
			authGroup.POST("/refresh", cfg.AuthHandler.Refresh)
		}

		// Records routes (Protected)
		records := api.Group("/records")
		records.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Viewer and above can list
			records.GET("", middleware.RoleMiddleware("VIEWER", "ANALYST", "ADMIN"), cfg.RecordHandler.ListRecords)
			
			// Analyst and Admin can create/delete
			records.POST("", middleware.RoleMiddleware("ANALYST", "ADMIN"), cfg.RecordHandler.CreateRecord)
			records.DELETE("/:id", middleware.RoleMiddleware("ADMIN"), cfg.RecordHandler.DeleteRecord)
		}

		// Dashboard routes (Protected)
		dashboard := api.Group("/dashboard")
		dashboard.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			dashboard.GET("/summary", middleware.RoleMiddleware("VIEWER", "ANALYST", "ADMIN"), cfg.DashboardHandler.GetSummary)
		}

		// Example protected route for testing roles
		protected := api.Group("/protected")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			protected.GET("/viewer", middleware.RoleMiddleware("VIEWER", "ANALYST", "ADMIN"), func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Hello Viewer/Analyst/Admin"})
			})
			protected.GET("/admin", middleware.RoleMiddleware("ADMIN"), func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Hello Admin Only"})
			})
		}
	}

	return r
}
