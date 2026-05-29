package db

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"clairvoyance/internal/auth"
	"clairvoyance/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func computeHash(recipient, title, credType, date string, skills []string) string {
	raw := fmt.Sprintf("%s|%s|%s|%s|%s",
		recipient,
		title,
		credType,
		date,
		strings.Join(skills, ","),
	)
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

func Seed(db *gorm.DB) error {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping seeder")
		return nil
	}

	log.Println("Seeding database...")

	// 1. Create Users
	hash, err := auth.HashPassword("demo123")
	if err != nil {
		return err
	}

	kuID := uuid.New()
	ku := &models.User{
		ID:            kuID,
		Email:         "issuer@kenyatta.edu",
		FullName:      "Kenyatta University",
		PasswordHash:  hash,
		Role:          "issuer",
		WalletAddress: "0x7f3a8b2c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a",
	}

	buildkeAdminID := uuid.New()
	buildkeAdmin := &models.User{
		ID:            buildkeAdminID,
		Email:         "admin@buildke.io",
		FullName:      "BuildKE Registrar",
		PasswordHash:  hash,
		Role:          "issuer",
		WalletAddress: "0x4d8ef21c9b0a7e6d5c4b3a2f1e0d9c8b7a6f5e4d",
	}

	aliceID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	alice := &models.User{
		ID:            aliceID,
		Email:         "alice@kenyatta.edu",
		FullName:      "Alice Wanjiku",
		PasswordHash:  hash,
		Role:          "holder",
		WalletAddress: "0x7f3a8b2c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a",
	}

	brianID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	brian := &models.User{
		ID:            brianID,
		Email:         "brian@buildke.io",
		FullName:      "Brian Omondi",
		PasswordHash:  hash,
		Role:          "holder",
		WalletAddress: "0x4d8ef21c9b0a7e6d5c4b3a2f1e0d9c8b7a6f5e4d",
	}

	carolID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	carol := &models.User{
		ID:            carolID,
		Email:         "carol@moringa.school",
		FullName:      "Carol Njeri",
		PasswordHash:  hash,
		Role:          "holder",
		WalletAddress: "0xa1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0",
	}

	users := []*models.User{ku, buildkeAdmin, alice, brian, carol}
	for _, u := range users {
		if err := db.Create(u).Error; err != nil {
			return err
		}
	}

	// 2. Create Team
	teamID := uuid.MustParse("77777777-7777-7777-7777-777777777777")
	team := &models.Team{
		ID:          teamID,
		LeaderID:    brianID,
		Name:        "BuildKE Solidity Devs",
		Slug:        "buildke",
		ProjectName: "ProofPass Platform",
	}
	if err := db.Create(team).Error; err != nil {
		return err
	}

	// 3. Create Team Members
	tmAlice := &models.TeamMember{
		ID:              uuid.New(),
		TeamID:          teamID,
		UserID:          aliceID,
		Role:            "Protocol Engineer",
		ContributionPct: 86,
	}
	tmBrian := &models.TeamMember{
		ID:              uuid.New(),
		TeamID:          teamID,
		UserID:          brianID,
		Role:            "Full-stack Builder",
		ContributionPct: 72,
	}
	tmCarol := &models.TeamMember{
		ID:              uuid.New(),
		TeamID:          teamID,
		UserID:          carolID,
		Role:            "Security Reviewer",
		ContributionPct: 64,
	}

	members := []*models.TeamMember{tmAlice, tmBrian, tmCarol}
	for _, m := range members {
		if err := db.Create(m).Error; err != nil {
			return err
		}
	}

	// 4. Create Sample Credentials
	credAlice := &models.Credential{
		ID:             uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		IssuerID:       kuID,
		RecipientRef:   "alice@kenyatta.edu",
		Title:          "BSc. Computer Science",
		CredentialType: "degree",
		Description:    "4-year programme, graduated with first class honours.",
		Skills:         []string{"Solidity", "Go", "Python"},
		IssueDate:      time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
		DataHash:       computeHash("alice@kenyatta.edu", "BSc. Computer Science", "degree", "2024-03-15", []string{"Solidity", "Go", "Python"}),
		TxHash:         "0x7f3a8b2c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a",
		BlockNumber:    4823901,
		CreatedAt:      time.Now(),
	}

	credBrian := &models.Credential{
		ID:             uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		IssuerID:       buildkeAdminID,
		RecipientRef:   "brian@buildke.io",
		Title:          "Solidity Developer — BuildKE 2025",
		CredentialType: "certificate",
		Description:    "Completed 12-week intensive Solidity & smart contract development bootcamp.",
		Skills:         []string{"Solidity", "Smart Contracts", "Hardhat"},
		IssueDate:      time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC),
		DataHash:       computeHash("brian@buildke.io", "Solidity Developer — BuildKE 2025", "certificate", "2025-05-20", []string{"Solidity", "Smart Contracts", "Hardhat"}),
		TxHash:         "0x4d8ef21c9b0a7e6d5c4b3a2f1e0d9c8b7a6f5e4d3c2b1a0f9e8d7c6b5a4f3e2d",
		BlockNumber:    4891234,
		CreatedAt:      time.Now(),
	}

	creds := []*models.Credential{credAlice, credBrian}
	for _, c := range creds {
		if err := db.Create(c).Error; err != nil {
			return err
		}
	}

	log.Println("Seeding complete successfully!")
	return nil
}
