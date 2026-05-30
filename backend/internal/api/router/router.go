package router

import (
	"log/slog"
	"net/http"

	"clairvoyance/internal/api/handlers"
	"clairvoyance/internal/api/middleware"
	"clairvoyance/internal/blockchain"
	"clairvoyance/internal/config"
	"clairvoyance/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "ProofPass API",
		})
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ── Blockchain client ──────────────────────────────────────────────────────
	// Use the real client when all three env vars are present.
	// Falls back to MockBlockchainClient during development so the server
	// starts without a live RPC connection.
	var chainClient blockchain.BlockchainClient
	if config.App.RPCUrl != "" && config.App.ContractAddress != "" && config.App.PrivateKey != "" {
		client, err := blockchain.NewClient(
			config.App.RPCUrl,
			config.App.ContractAddress,
			config.App.PrivateKey,
		)
		if err != nil {
			slog.Warn("failed to init blockchain client — falling back to mock",
				"error", err,
			)
			chainClient = &blockchain.MockBlockchainClient{}
		} else {
			slog.Info("blockchain client connected", "rpc", config.App.RPCUrl)
			chainClient = client
		}
	} else {
		slog.Warn("RPC_URL / CONTRACT_ADDRESS / PRIVATE_KEY not set — using MockBlockchainClient")
		chainClient = &blockchain.MockBlockchainClient{}
	}

	// ── Services ───────────────────────────────────────────────────────────────
	credService := services.NewCredentialService(db, chainClient)
	teamService := services.NewTeamService(db)

	// ── Handlers ───────────────────────────────────────────────────────────────
	authHandler := handlers.NewAuthHandler(db)
	credHandler := handlers.NewCredentialHandler(credService)
	teamHandler := handlers.NewTeamHandler(teamService, credService)

	// ── Public routes ──────────────────────────────────────────────────────────
	api := r.Group("/api")
	{
		// Auth
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/refresh", authHandler.Refresh)

		// Credentials — verify is always public (no login, no MetaMask)
		api.GET("/credentials/verify/:id", credHandler.Verify)

		// Teams — read-only team data is public
		api.GET("/teams/:id", teamHandler.GetTeam)
		api.GET("/teams/:id/credentials", teamHandler.GetTeamCredentials)
	}

	// ── Protected routes ───────────────────────────────────────────────────────
	protected := api.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		// Auth
		protected.POST("/auth/logout", authHandler.Logout)
		protected.GET("/auth/profile", authHandler.Profile)

		// Credentials
		protected.GET("/credentials/mine", credHandler.GetMine)
		protected.POST("/credentials/issue", credHandler.Issue)
		protected.GET("/credentials/issued", credHandler.GetIssued)
	}

	return r
}
