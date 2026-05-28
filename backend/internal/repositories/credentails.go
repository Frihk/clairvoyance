package repository

import (
	"errors"

	"clairvoyance/internal/models"
	"clairvoyance/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CredentialRepository struct {
	db *gorm.DB
}

func NewCredentialRepository(db *gorm.DB) *CredentialRepository {
	return &CredentialRepository{db: db}
}

// Create inserts a new credential record
func (cr *CredentialRepository) Create(credential *models.Credential) error {
	result := cr.db.Create(credential)
	return result.Error
}

// GetByID retrieves a credential by UUID
func (cr *CredentialRepository) GetByID(id uuid.UUID) (*models.Credential, error) {
	var credential models.Credential
	result := cr.db.Where("id = ?", id).First(&credential)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, result.Error
	}
	return &credential, nil
}

// GetByRecipientRef retrieves all credentials for a given recipient (email or wallet)
func (cr *CredentialRepository) GetByRecipientRef(recipientRef string) ([]models.Credential, error) {
	var credentials []models.Credential
	result := cr.db.Where("recipient_ref = ?", recipientRef).
		Order("created_at DESC").
		Find(&credentials)
	return credentials, result.Error
}

// GetByIssuerID retrieves all credentials issued by a specific issuer
func (cr *CredentialRepository) GetByIssuerID(issuerID uuid.UUID) ([]models.Credential, error) {
	var credentials []models.Credential
	result := cr.db.Where("issuer_id = ?", issuerID).
		Order("created_at DESC").
		Find(&credentials)
	return credentials, result.Error
}

// GetByRecipientRefs retrieves credentials for multiple recipients (used for team credentials)
func (cr *CredentialRepository) GetByRecipientRefs(recipientRefs []string) ([]models.Credential, error) {
	if len(recipientRefs) == 0 {
		return []models.Credential{}, nil
	}
	var credentials []models.Credential
	result := cr.db.Where("recipient_ref IN ?", recipientRefs).
		Order("created_at DESC").
		Find(&credentials)
	return credentials, result.Error
}

// Update updates an existing credential
func (cr *CredentialRepository) Update(credential *models.Credential) error {
	result := cr.db.Save(credential)
	return result.Error
}

// GetByDataHash retrieves a credential by its unique data hash
func (cr *CredentialRepository) GetByDataHash(dataHash string) (*models.Credential, error) {
	var credential models.Credential
	result := cr.db.Where("data_hash = ?", dataHash).First(&credential)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, result.Error
	}
	return &credential, nil
}