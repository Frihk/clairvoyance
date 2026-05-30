package app

import (
	"log"
	"net/http"

	"clairvoyance/internal/api/router"
	"clairvoyance/internal/auth"
	"clairvoyance/internal/config"
	"clairvoyance/internal/models"
	"clairvoyance/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewHandler() (http.Handler, error) {
	config.Load()
	utils.InitLogger(config.App.AppEnv)

	db, err := gorm.Open(postgres.Open(config.App.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Credential{},
		&models.Team{},
		&models.TeamMember{},
	); err != nil {
		return nil, err
	}

	if err := seedDemoIssuer(db); err != nil {
		return nil, err
	}

	return router.SetupRouter(db), nil
}

func MustHandler() http.Handler {
	handler, err := NewHandler()
	if err != nil {
		log.Fatal(err)
	}

	return handler
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
