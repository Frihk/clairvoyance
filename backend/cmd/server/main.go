package main

import (
	"log"

	"clairvoyance/internal/api/router"
	"clairvoyance/internal/config"
	"clairvoyance/internal/models"
	"clairvoyance/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config.Load()
	utils.InitLogger(config.App.AppEnv)

	// Database connection
	db, err := gorm.Open(postgres.Open(config.App.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate all models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Credential{},
		&models.Team{},
		&models.TeamMember{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := router.SetupRouter(db)
	if err := r.Run(":" + config.App.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}