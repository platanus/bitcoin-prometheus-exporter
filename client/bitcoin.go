package client

import (
	"fmt"
	"log"

	"github.com/btcsuite/btcd/rpcclient"
)

// BitcoinClient allows you to fetch Bitcoin node metrics from rpc.
type BitcoinClient struct {
	rpcclient *rpcclient.Client
}

// NodeStats represents node metrics.
type NodeStats struct {
	BlockCount      int64
	ConnectionCount int64
	Difficulty      float64
}

// NewBitcoinClient creates an BitcoinClient.
func NewBitcoinClient(rpcclient *rpcclient.Client) (*BitcoinClient, error) {
	client := &BitcoinClient{
		rpcclient: rpcclient,
	}

	if _, err := client.GetStats(); err != nil {
		return nil, fmt.Errorf("Failed to create BitcoinClient: %v", err)
	}

	return client, nil
}

// GetStats fetches the node metrics.
func (client *BitcoinClient) GetStats() (*NodeStats, error) {

	var stats NodeStats

	blockCount, err := client.rpcclient.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	stats.BlockCount = blockCount

	connectionCount, err := client.rpcclient.GetConnectionCount()
	if err != nil {
		log.Fatal(err)
	}
	stats.ConnectionCount = connectionCount

	difficulty, err := client.rpcclient.GetDifficulty()
	if err != nil {
		log.Fatal(err)
	}
	stats.Difficulty = difficulty

	return &stats, nil
}
