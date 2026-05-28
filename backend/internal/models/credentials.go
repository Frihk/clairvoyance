package models

import (
	"time"

	"github.com/google/uuid"
)

type Credential struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	IssuerID       uuid.UUID `gorm:"type:uuid;not null"   json:"issuer_id"`
	RecipientRef   string    `gorm:"not null"             json:"recipient_ref"`
	Title          string    `gorm:"not null"             json:"title"`
	CredentialType string    `gorm:"not null"             json:"credential_type"`
	Description    string    `json:"description"`
	Skills         []string  `gorm:"serializer:json"      json:"skills"`
	IssueDate      time.Time `gorm:"not null"             json:"issue_date"`
	DataHash       string    `gorm:"uniqueIndex;not null" json:"-"`
	TxHash         string    `gorm:"not null"             json:"tx_hash"`
	BlockNumber    int64     `gorm:"not null"             json:"block_number"`
	CreatedAt      time.Time `json:"created_at"`
}