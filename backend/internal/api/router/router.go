package router

import (
	"clairvoyance/internal/api/handlers"
	"clairvoyance/internal/api/middleware"
	"clairvoyance/internal/blockchain"
	"clairvoyance/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Initialize blockchain client (replace mock with real implementation)
	chainClient := &blockchain.MockBlockchainClient{}

	// Initialize services
	credService := services.NewCredentialService(db, chainClient)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	credHandler := handlers.NewCredentialHandler(credService)

	// Public routes (no auth required)
	api := r.Group("/api")
	{
		// Auth endpoints
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/refresh", authHandler.Refresh)

		// Public verify endpoint - CRITICAL for MVP
		api.GET("/credentials/verify/:id", credHandler.Verify)
	}

	// Protected routes (auth required)
	protected := api.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		// Auth
		protected.POST("/auth/logout", authHandler.Logout)
		protected.GET("/auth/profile", authHandler.Profile)

		// Credentials
		protected.GET("/credentials/mine", credHandler.GetMine)
		protected.POST("/credentials/issue", credHandler.Issue) // Requires issuer role
		protected.GET("/credentials/issued", credHandler.GetIssued)
	}

	return r
}
