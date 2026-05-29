package repository

import (
	"errors"

	"clairvoyance/internal/models"
	"clairvoyance/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

// GetByID returns a team with its members and each member's user record preloaded.
func (tr *TeamRepository) GetByID(id uuid.UUID) (*models.Team, error) {
	var team models.Team
	result := tr.db.
		Preload("Members").
		Preload("Members.User").
		Where("id = ?", id).
		First(&team)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, result.Error
	}
	return &team, nil
}

// GetBySlug returns a team by its URL-safe slug.
func (tr *TeamRepository) GetBySlug(slug string) (*models.Team, error) {
	var team models.Team
	result := tr.db.
		Preload("Members").
		Preload("Members.User").
		Where("slug = ?", slug).
		First(&team)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, result.Error
	}
	return &team, nil
}

// GetMembersByTeamID returns all TeamMember rows for a given team,
// with each member's User preloaded.
func (tr *TeamRepository) GetMembersByTeamID(teamID uuid.UUID) ([]models.TeamMember, error) {
	var members []models.TeamMember
	result := tr.db.
		Preload("User").
		Where("team_id = ?", teamID).
		Find(&members)
	return members, result.Error
}
