package main

import (
	"log"

	"clairvoyance/internal/api/router"
	"clairvoyance/internal/auth"
	"clairvoyance/internal/config"
	"clairvoyance/internal/models"
	"clairvoyance/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	if err := seedDemoIssuer(db); err != nil {
		log.Fatal("Failed to seed demo issuer:", err)
	}

	r := router.SetupRouter(db)
	if err := r.Run(":" + config.App.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func seedDemoIssuer(db *gorm.DB) error {
	if config.App.IssuerEmail == "" || config.App.IssuerPass == "" {
		return nil
	}

	passwordHash, err := auth.HashPassword(config.App.IssuerPass)
	if err != nil {
		return err
	}

	user := models.User{
		Email:        config.App.IssuerEmail,
		PasswordHash: passwordHash,
		FullName:     "Demo Issuer",
		Role:         "issuer",
	}

	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "email"}},
		DoUpdates: clause.Assignments(map[string]any{
			"password_hash": passwordHash,
			"role":          "issuer",
		}),
	}).Create(&user).Error
}
