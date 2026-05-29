package blockchain

import (
	"context"
	"log/slog"
)

// BlockchainClient is the interface the service layer uses to talk to the
// ProofPass smart contract. The real implementation is *Client in contract.go.
// MockBlockchainClient below is used by the router until a live RPC URL and
// private key are available.
type BlockchainClient interface {
	IssueCredential(ctx context.Context, credentialID, dataHash string) (*IssueResult, error)
	VerifyCredential(ctx context.Context, credentialID string) (*VerifyResult, error)
}

// MockBlockchainClient satisfies BlockchainClient without touching the chain.
// Replace it in router.go with a real *Client once the .env has RPC_URL,
// CONTRACT_ADDRESS, and PRIVATE_KEY set.
type MockBlockchainClient struct{}

func (m *MockBlockchainClient) IssueCredential(_ context.Context, credentialID, dataHash string) (*IssueResult, error) {
	slog.Warn("MockBlockchainClient.IssueCredential — blockchain not wired yet",
		"credential_id", credentialID,
		"data_hash", dataHash,
	)
	// Returns a clearly fake tx hash so the verify flow returns TAMPERED
	// rather than crashing during development. Replace with real client ASAP.
	return &IssueResult{TxHash: "0x0000000000000000000000000000000000000000000000000000000000000000"}, nil
}

func (m *MockBlockchainClient) VerifyCredential(_ context.Context, credentialID string) (*VerifyResult, error) {
	slog.Warn("MockBlockchainClient.VerifyCredential — blockchain not wired yet",
		"credential_id", credentialID,
	)
	return &VerifyResult{Found: false}, nil
}
