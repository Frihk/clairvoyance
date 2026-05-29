package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Credential struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"         json:"id"`
	IssuerID       uuid.UUID `gorm:"type:uuid;not null"           json:"issuer_id"`
	RecipientRef   string    `gorm:"not null"                     json:"recipient_ref"`
	Title          string    `gorm:"not null"                     json:"title"`
	CredentialType string    `gorm:"not null"                     json:"credential_type"`
	Description    string    `                                    json:"description"`
	Skills         []string  `gorm:"serializer:json"              json:"skills"`
	IssueDate      time.Time `gorm:"not null"                     json:"issue_date"`
	DataHash       string    `gorm:"uniqueIndex;not null;size:64" json:"-"`
	TxHash         string    `gorm:"not null;size:66"             json:"tx_hash"`
	BlockNumber    int64     `gorm:"not null"                     json:"block_number"`
	CreatedAt      time.Time `                                    json:"created_at"`

	// Loaded via Preload
	Issuer *User `gorm:"foreignKey:IssuerID" json:"issuer,omitempty"`
}

func (c *Credential) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
