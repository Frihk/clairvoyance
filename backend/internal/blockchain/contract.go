package blockchain

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

const defaultTxTimeout = 60 * time.Second

type IssueResult struct {
	TxHash string
}

type VerifyResult struct {
	DataHash string
	Issuer   string
	IssuedAt time.Time
	Found    bool
}

func (c *Client) IssueCredential(ctx context.Context, credentialID, dataHash string) (*IssueResult, error) {
	if c == nil {
		return nil, fmt.Errorf("blockchain client is nil")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	auth := *c.auth
	auth.Context = ctx

	tx, err := c.contract.IssueCredential(&auth, credentialID, dataHash)
	if err != nil {
		return nil, fmt.Errorf("issueCredential tx: %w", err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, defaultTxTimeout)
	defer cancel()

	receipt, err := bind.WaitMined(waitCtx, c.eth, tx)
	if err != nil {
		return nil, fmt.Errorf("wait mined: %w", err)
	}

	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction reverted: %s", tx.Hash().Hex())
	}

	return &IssueResult{TxHash: tx.Hash().Hex()}, nil
}

func (c *Client) VerifyCredential(ctx context.Context, credentialID string) (*VerifyResult, error) {
	if c == nil {
		return nil, fmt.Errorf("blockchain client is nil")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	dataHash, issuerAddr, issuedAtBig, err := c.contract.VerifyCredential(
		&bind.CallOpts{Context: ctx},
		credentialID,
	)
	if err != nil {
		return nil, fmt.Errorf("verifyCredential call: %w", err)
	}

	if issuedAtBig.Sign() == 0 {
		return &VerifyResult{Found: false}, nil
	}

	return &VerifyResult{
		DataHash: dataHash,
		Issuer:   issuerAddr.Hex(),
		IssuedAt: time.Unix(issuedAtBig.Int64(), 0),
		Found:    true,
	}, nil
}
