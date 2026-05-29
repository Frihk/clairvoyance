package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"clairvoyance/internal/blockchain"
	"clairvoyance/internal/models"
	repository "clairvoyance/internal/repositories"
	"clairvoyance/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CredentialService struct {
	credRepo *repository.CredentialRepository
	userRepo *repository.UserRepository
	chain    blockchain.BlockchainClient
}

func NewCredentialService(db *gorm.DB, chain blockchain.BlockchainClient) *CredentialService {
	return &CredentialService{
		credRepo: repository.NewCredentialRepository(db),
		userRepo: repository.NewUserRepository(db),
		chain:    chain,
	}
}

// IssueRequest represents the data needed to issue a new credential
type IssueRequest struct {
	Recipient      string    `json:"recipient" binding:"required"` // email or wallet address
	CredentialType string    `json:"credential_type" binding:"required"`
	Title          string    `json:"title" binding:"required"`
	Description    string    `json:"description"`
	IssueDate      time.Time `json:"issue_date" binding:"required"`
	Skills         []string  `json:"skills"`
}

// IssueResponse represents the response after issuing a credential
type IssueResponse struct {
	ID          uuid.UUID `json:"id"`
	VerifyURL   string    `json:"verify_url"`
	TxHash      string    `json:"tx_hash"`
	BlockNumber int64     `json:"block_number"`
	DataHash    string    `json:"data_hash"`
}

// VerifyResponse represents the verification result
type VerifyResponse struct {
	Status       string                `json:"status"` // VERIFIED, TAMPERED, NOT_FOUND, ERROR
	Message      string                `json:"message,omitempty"`
	Credential   *models.Credential    `json:"credential,omitempty"`
	BlockNumber  int64                 `json:"block_number,omitempty"`
	TxHash       string                `json:"tx_hash,omitempty"`
	Network      string                `json:"network,omitempty"`
	ComputedHash string                `json:"computed_hash,omitempty"`
	OnChainHash  string                `json:"onchain_hash,omitempty"`
}

// IssueCredential creates a new credential, writes to blockchain, and stores in DB
func (cs *CredentialService) IssueCredential(issuerID uuid.UUID, req IssueRequest) (*IssueResponse, error) {
	// Verify issuer exists
	issuer, err := cs.userRepo.GetUserByID(issuerID)
	if err != nil {
		return nil, fmt.Errorf("issuer not found: %w", err)
	}
	if issuer.Role != "issuer" && issuer.Role != "leader" {
		return nil, utils.ErrForbidden
	}

	// Create credential object
	credentialID := uuid.New()
	credential := &models.Credential{
		ID:             credentialID,
		IssuerID:       issuerID,
		RecipientRef:   req.Recipient,
		Title:          req.Title,
		CredentialType: req.CredentialType,
		Description:    req.Description,
		Skills:         req.Skills,
		IssueDate:      req.IssueDate,
		CreatedAt:      time.Now(),
	}

	// Compute canonical hash
	dataHash := computeCanonicalHash(credential)
	credential.DataHash = dataHash

	// Write to blockchain
	credentialIDStr := credentialID.String()
	result, err := cs.chain.IssueCredential(context.Background(), credentialIDStr, dataHash)
	if err != nil {
		slog.Error("blockchain write failed", "error", err, "credential_id", credentialIDStr)
		return nil, fmt.Errorf("blockchain write failed: %w", err)
	}
	credential.TxHash = result.TxHash
	credential.BlockNumber = result.BlockNumber

	// Store in database
	if err := cs.credRepo.Create(credential); err != nil {
		slog.Error("database save failed", "error", err, "credential_id", credentialID)
		// Note: Blockchain already has the credential, but DB failed.
		// In production, you'd want to handle this inconsistency.
		return nil, fmt.Errorf("database save failed: %w", err)
	}

	slog.Info("credential issued", 
		"credential_id", credentialID, 
		"issuer_id", issuerID, 
		"recipient", req.Recipient,
		"tx_hash", result.TxHash)

	// Build verify URL (configure base URL from env)
	verifyURL := fmt.Sprintf("/verify/%s", credentialIDStr)

	return &IssueResponse{
		ID:          credentialID,
		VerifyURL:   verifyURL,
		TxHash:      result.TxHash,
		BlockNumber: result.BlockNumber,
		DataHash:    dataHash,
	}, nil
}

// VerifyCredential verifies a credential by comparing stored data with on-chain record
func (cs *CredentialService) VerifyCredential(credentialID uuid.UUID) (*VerifyResponse, error) {
	// Load from database
	credential, err := cs.credRepo.GetByID(credentialID)
	if err != nil {
		if err == utils.ErrNotFound {
			return &VerifyResponse{
				Status:  "NOT_FOUND",
				Message: "Credential not found",
			}, nil
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Recompute hash from stored data
	recomputedHash := computeCanonicalHash(credential)

	// Fetch on-chain hash
	verifyResult, err := cs.chain.VerifyCredential(context.Background(), credentialID.String())
	if err != nil {
		slog.Error("blockchain read failed", "error", err, "credential_id", credentialID)
		return &VerifyResponse{
			Status:  "ERROR",
			Message: "Failed to read from blockchain",
		}, nil
	}

	// Compare
	if !verifyResult.Found {
		return &VerifyResponse{
			Status:  "NOT_FOUND",
			Message: "Credential not found on chain",
		}, nil
	}

	if recomputedHash == verifyResult.DataHash {
		return &VerifyResponse{
			Status:      "VERIFIED",
			Credential:  credential,
			BlockNumber: credential.BlockNumber,
			TxHash:      credential.TxHash,
			Network:     "Polygon Amoy",
		}, nil
	}

	// Tampered
	return &VerifyResponse{
		Status:       "TAMPERED",
		Message:      "Credential data does not match blockchain record",
		Credential:   credential,
		ComputedHash: recomputedHash,
		OnChainHash:  verifyResult.DataHash,
	}, nil
}

// GetUserCredentials returns all credentials owned by a user (by email or wallet)
func (cs *CredentialService) GetUserCredentials(userID uuid.UUID) ([]models.Credential, error) {
	user, err := cs.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	var allCredentials []models.Credential
	recipientRefs := []string{}

	if user.Email != "" {
		recipientRefs = append(recipientRefs, user.Email)
	}
	if user.WalletAddress != "" {
		recipientRefs = append(recipientRefs, user.WalletAddress)
	}

	for _, ref := range recipientRefs {
		creds, err := cs.credRepo.GetByRecipientRef(ref)
		if err != nil {
			continue
		}
		allCredentials = append(allCredentials, creds...)
	}

	// Deduplicate by ID
	seen := make(map[uuid.UUID]bool)
	unique := []models.Credential{}
	for _, cred := range allCredentials {
		if !seen[cred.ID] {
			seen[cred.ID] = true
			unique = append(unique, cred)
		}
	}

	return unique, nil
}

// GetIssuerCredentials returns all credentials issued by a specific issuer
func (cs *CredentialService) GetIssuerCredentials(issuerID uuid.UUID) ([]models.Credential, error) {
	return cs.credRepo.GetByIssuerID(issuerID)
}

// GetTeamCredentials returns credentials for all team members
func (cs *CredentialService) GetTeamCredentials(teamID uuid.UUID) ([]models.Credential, error) {
	// This requires a team repository - implement when needed
	// For now, return empty
	return []models.Credential{}, nil
}

// computeCanonicalHash creates a deterministic hash from credential data
// Order must be consistent across issue and verify operations
func computeCanonicalHash(cred *models.Credential) string {
	raw := fmt.Sprintf("%s|%s|%s|%s|%s",
		cred.RecipientRef,
		cred.Title,
		cred.CredentialType,
		cred.IssueDate.Format("2006-01-02"),
		strings.Join(cred.Skills, ","),
	)
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}
