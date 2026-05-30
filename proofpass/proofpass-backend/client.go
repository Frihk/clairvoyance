
package blockchain

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	eth      *ethclient.Client
	contract *ProofPass
	auth     *bind.TransactOpts
}

func NewClient(rpcURL, contractAddr, privateKey string) (*Client, error) {
	eth, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial rpc: %w", err)
	}

	addr := common.HexToAddress(contractAddr)

	contract, err := NewProofPass(addr, eth)
	if err != nil {
		return nil, fmt.Errorf("load contract: %w", err)
	}

	privateKey = strings.TrimPrefix(privateKey, "0x")
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	chainID, err := eth.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("read chain id: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("build transactor: %w", err)
	}

	return &Client{eth: eth, contract: contract, auth: auth}, nil
}
