package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"   json:"id"`
	LeaderID    uuid.UUID `gorm:"type:uuid;not null"     json:"leader_id"`
	Name        string    `gorm:"not null"               json:"name"`
	Slug        string    `gorm:"uniqueIndex;not null"   json:"slug"`
	ProjectName string    `gorm:"not null"               json:"project_name"`
	CreatedAt   time.Time `                              json:"created_at"`

	// Loaded via Preload — not stored in the teams table
	Members []TeamMember `gorm:"foreignKey:TeamID" json:"members,omitempty"`
}

type TeamMember struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"      json:"id"`
	TeamID          uuid.UUID `gorm:"type:uuid;not null;index"  json:"team_id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null"        json:"user_id"`
	Role            string    `gorm:"not null"                  json:"role"`
	ContributionPct float64   `gorm:"not null"                  json:"contribution_pct"`

	// Loaded via Preload
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (t *Team) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (tm *TeamMember) BeforeCreate(tx *gorm.DB) error {
	if tm.ID == uuid.Nil {
		tm.ID = uuid.New()
	}
	return nil
}
