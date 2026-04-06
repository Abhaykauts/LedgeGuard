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
	UserHandler      *UserHandler
	JWTSecret        string
}

// NewRouter sets up the API routes
func NewRouter(cfg RouterConfig) *gin.Engine {
	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

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

		// User routes (Admin Only)
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			users.GET("", middleware.RoleMiddleware("ADMIN"), cfg.UserHandler.ListUsers)
			users.POST("", middleware.RoleMiddleware("ADMIN"), cfg.UserHandler.CreateUser)
			users.PUT("/:id", middleware.RoleMiddleware("ADMIN"), cfg.UserHandler.UpdateUser)
			users.DELETE("/:id", middleware.RoleMiddleware("ADMIN"), cfg.UserHandler.DeleteUser)
		}

		// Records routes (Protected)
		records := api.Group("/records")
		records.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Viewer and above can list
			records.GET("", middleware.RoleMiddleware("VIEWER", "ANALYST", "ADMIN"), cfg.RecordHandler.ListRecords)

			// Admin can create/update/delete
			records.POST("", middleware.RoleMiddleware("ADMIN"), cfg.RecordHandler.CreateRecord)
			records.PUT("/:id", middleware.RoleMiddleware("ADMIN"), cfg.RecordHandler.UpdateRecord)
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
