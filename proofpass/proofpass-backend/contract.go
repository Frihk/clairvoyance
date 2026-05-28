
package blockchain

import (

	"context"

	"fmt"

	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

)

type IssueResult struct {

	TxHash string

}

type VerifyResult struct {

	DataHash string

	Issuer   string

	IssuedAt time.Time

	Found    bool

}

func (c *Client) IssueCredential(credentialID, dataHash string) (*IssueResult, error) {

	tx, err := c.contract.IssueCredential(c.auth, credentialID, dataHash)

	if err != nil {

		return nil, fmt.Errorf("issueCredential tx: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()

	receipt, err := bind.WaitMined(ctx, c.eth, tx)

	if err != nil {

		return nil, fmt.Errorf("wait mined: %w", err)

	}

	if receipt.Status == 0 {

		return nil, fmt.Errorf("transaction reverted: %s", tx.Hash().Hex())

	}

	return &IssueResult{TxHash: tx.Hash().Hex()}, nil

}

func (c *Client) VerifyCredential(credentialID string) (*VerifyResult, error) {

	dataHash, issuerAddr, issuedAtBig, err := c.contract.VerifyCredential(

		&bind.CallOpts{Context: context.Background()},

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

