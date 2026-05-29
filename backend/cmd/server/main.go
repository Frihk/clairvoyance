package main

import (
	"log"

	"clairvoyance/internal/api/router"
	"clairvoyance/internal/config"
	"clairvoyance/internal/db"
	"clairvoyance/internal/models"
	"clairvoyance/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config.Load()
	utils.InitLogger(config.App.AppEnv)

	// Database connection
	dbConn, err := gorm.Open(postgres.Open(config.App.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate all models
	if err := dbConn.AutoMigrate(
		&models.User{},
		&models.Credential{},
		&models.Team{},
		&models.TeamMember{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed demo data
	if err := db.Seed(dbConn); err != nil {
		log.Println("Warning: Failed to seed database:", err)
	}

	r := router.SetupRouter(dbConn)
	if err := r.Run(":" + config.App.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
