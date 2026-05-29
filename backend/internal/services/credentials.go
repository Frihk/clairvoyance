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
	"clairvoyance/internal/config"
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
	Recipient      string    `json:"recipient"       binding:"required"` // email or wallet address
	CredentialType string    `json:"credential_type" binding:"required"`
	Title          string    `json:"title"           binding:"required"`
	Description    string    `json:"description"`
	IssueDate      time.Time `json:"issue_date"      binding:"required"`
	Skills         []string  `json:"skills"`
}

// IssueResponse is returned after a credential is successfully issued.
type IssueResponse struct {
	ID          uuid.UUID `json:"id"`
	VerifyURL   string    `json:"verify_url"`
	TxHash      string    `json:"tx_hash"`
	BlockNumber int64     `json:"block_number"`
	DataHash    string    `json:"data_hash"`
}

// VerifyResponse is the structured result of a credential verification request.
type VerifyResponse struct {
	Status       string             `json:"status"` // VERIFIED | TAMPERED | NOT_FOUND | ERROR
	Message      string             `json:"message,omitempty"`
	Credential   *models.Credential `json:"credential,omitempty"`
	BlockNumber  int64              `json:"block_number,omitempty"`
	TxHash       string             `json:"tx_hash,omitempty"`
	Network      string             `json:"network,omitempty"`
	ComputedHash string             `json:"computed_hash,omitempty"` // only present on TAMPERED
	OnChainHash  string             `json:"onchain_hash,omitempty"`  // only present on TAMPERED
}

// IssueCredential hashes the credential payload, writes it to the ProofPass
// contract on Polygon, then persists the full record to the database.
func (cs *CredentialService) IssueCredential(issuerID uuid.UUID, req IssueRequest) (*IssueResponse, error) {
	// Verify issuer exists and has the right role
	issuer, err := cs.userRepo.GetUserByID(issuerID)
	if err != nil {
		return nil, fmt.Errorf("issuer not found: %w", err)
	}
	if issuer.Role != "issuer" && issuer.Role != "leader" {
		return nil, utils.ErrForbidden
	}

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

	// Compute canonical hash before writing to chain
	dataHash := computeCanonicalHash(credential)
	credential.DataHash = dataHash

	// Write hash to Polygon
	credentialIDStr := credentialID.String()
	result, err := cs.chain.IssueCredential(context.Background(), credentialIDStr, dataHash)
	if err != nil {
		slog.Error("blockchain write failed", "error", err, "credential_id", credentialIDStr)
		return nil, fmt.Errorf("blockchain write failed: %w", err)
	}
	credential.TxHash = result.TxHash
	credential.BlockNumber = result.BlockNumber

	// Persist to DB
	if err := cs.credRepo.Create(credential); err != nil {
		slog.Error("database save failed", "error", err, "credential_id", credentialID)
		return nil, fmt.Errorf("database save failed: %w", err)
	}

	slog.Info("credential issued",
		"credential_id", credentialID,
		"issuer_id", issuerID,
		"recipient", req.Recipient,
		"tx_hash", result.TxHash,
	)

	base := config.App.AppBaseURL
	verifyURL := fmt.Sprintf("%s/api/credentials/verify/%s", base, credentialIDStr)

	return &IssueResponse{
		ID:          credentialID,
		VerifyURL:   verifyURL,
		TxHash:      result.TxHash,
		BlockNumber: result.BlockNumber,
		DataHash:    dataHash,
	}, nil
}

// VerifyCredential re-hashes the stored credential data and compares it
// against the on-chain record. Returns VERIFIED, TAMPERED, NOT_FOUND, or ERROR.
func (cs *CredentialService) VerifyCredential(credentialID uuid.UUID) (*VerifyResponse, error) {
	credential, err := cs.credRepo.GetByID(credentialID)
	if err != nil {
		if err == utils.ErrNotFound {
			return &VerifyResponse{Status: "NOT_FOUND", Message: "Credential not found"}, nil
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	recomputedHash := computeCanonicalHash(credential)

	verifyResult, err := cs.chain.VerifyCredential(context.Background(), credentialID.String())
	if err != nil {
		slog.Error("blockchain read failed", "error", err, "credential_id", credentialID)
		return &VerifyResponse{Status: "ERROR", Message: "Failed to read from blockchain"}, nil
	}

	if !verifyResult.Found {
		return &VerifyResponse{Status: "NOT_FOUND", Message: "Credential not found on chain"}, nil
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

	slog.Warn("hash mismatch",
		"credential_id", credentialID,
		"recomputed", recomputedHash,
		"on_chain", verifyResult.DataHash,
	)
	return &VerifyResponse{
		Status:       "TAMPERED",
		Message:      "Credential data does not match blockchain record",
		Credential:   credential,
		ComputedHash: recomputedHash,
		OnChainHash:  verifyResult.DataHash,
	}, nil
}

// GetUserCredentials returns all credentials where recipient_ref matches
// the user's email or wallet address, deduplicating by ID.
func (cs *CredentialService) GetUserCredentials(userID uuid.UUID) ([]models.Credential, error) {
	user, err := cs.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	refs := []string{}
	if user.Email != "" {
		refs = append(refs, user.Email)
	}
	if user.WalletAddress != "" {
		refs = append(refs, user.WalletAddress)
	}

	return cs.credRepo.GetByRecipientRefs(refs)
}

// GetIssuerCredentials returns all credentials issued by the given user.
func (cs *CredentialService) GetIssuerCredentials(issuerID uuid.UUID) ([]models.Credential, error) {
	return cs.credRepo.GetByIssuerID(issuerID)
}

// GetCredentialsByRefs returns all credentials for a list of recipient refs
// (emails or wallet addresses). Used by the team credentials endpoint.
func (cs *CredentialService) GetCredentialsByRefs(refs []string) ([]models.Credential, error) {
	return cs.credRepo.GetByRecipientRefs(refs)
}

// computeCanonicalHash produces a deterministic SHA-256 digest of a
// credential's core fields. Field order is fixed — never reorder after deploy
// or all existing on-chain hashes will stop verifying.
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
