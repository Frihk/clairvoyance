package services

import (
	"clairvoyance/internal/models"
	repository "clairvoyance/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamService struct {
	teamRepo *repository.TeamRepository
}

func NewTeamService(db *gorm.DB) *TeamService {
	return &TeamService{
		teamRepo: repository.NewTeamRepository(db),
	}
}

// GetTeamByID returns the team with members preloaded.
func (ts *TeamService) GetTeamByID(id uuid.UUID) (*models.Team, error) {
	return ts.teamRepo.GetByID(id)
}

// GetTeamBySlug returns the team by URL slug with members preloaded.
func (ts *TeamService) GetTeamBySlug(slug string) (*models.Team, error) {
	return ts.teamRepo.GetBySlug(slug)
}

// GetMemberRefs returns all email addresses and wallet addresses belonging to
// a team's members. Used by CredentialService to look up team credentials.
func (ts *TeamService) GetMemberRefs(teamID uuid.UUID) ([]string, error) {
	members, err := ts.teamRepo.GetMembersByTeamID(teamID)
	if err != nil {
		return nil, err
	}

	refs := make([]string, 0, len(members)*2)
	for _, m := range members {
		if m.User == nil {
			continue
		}
		if m.User.Email != "" {
			refs = append(refs, m.User.Email)
		}
		if m.User.WalletAddress != "" {
			refs = append(refs, m.User.WalletAddress)
		}
	}
	return refs, nil
}
