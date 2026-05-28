package blockchain

import "time"

// BlockchainClient defines the interface for blockchain operations
// Person 01 will implement this with go-ethereum
type BlockchainClient interface {
	// IssueCredential writes a credential hash to the blockchain
	// Returns: transaction hash, block number, error
	IssueCredential(credentialID, dataHash string) (txHash string, blockNumber int64, err error)
	
	// VerifyCredential reads a credential hash from the blockchain
	// Returns: data hash, issuer address, timestamp, error
	VerifyCredential(credentialID string) (dataHash string, issuerAddress string, issuedAt uint64, err error)
	
	// AddIssuer authorizes a new issuer (only callable by contract owner)
	AddIssuer(issuerAddress string) error
	
	// IsIssuer checks if an address is an authorized issuer
	IsIssuer(address string) (bool, error)
}

// Mock implementation for development/testing
type MockBlockchainClient struct{}

func (m *MockBlockchainClient) IssueCredential(credentialID, dataHash string) (string, int64, error) {
	// Mock implementation - replace with real go-ethereum calls
	return "0xmocktxhash1234567890", 12345678, nil
}

func (m *MockBlockchainClient) VerifyCredential(credentialID string) (string, string, uint64, error) {
	// Mock implementation - replace with real go-ethereum calls
	return "mockhash", "0xissueraddress", uint64(time.Now().Unix()), nil
}

func (m *MockBlockchainClient) AddIssuer(issuerAddress string) error {
	return nil
}

func (m *MockBlockchainClient) IsIssuer(address string) (bool, error) {
	return true, nil
}