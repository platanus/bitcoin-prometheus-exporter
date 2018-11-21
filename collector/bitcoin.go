package collector

import (
	"log"
	"sync"

	"github.com/platanus/bitcoin-prometheus-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

// BitcoinCollector collects node metrics. It implements prometheus.Collector interface.
type BitcoinCollector struct {
	bitcoinClient *client.BitcoinClient
	metrics       map[string]*prometheus.Desc
	mutex         sync.Mutex
}

// NewBitcoinCollector creates an BitcoinCollector.
func NewBitcoinCollector(bitcoinClient *client.BitcoinClient, namespace string) *BitcoinCollector {
	return &BitcoinCollector{
		bitcoinClient: bitcoinClient,
		metrics: map[string]*prometheus.Desc{
			"block_total":      newGlobalMetric(namespace, "block_total", "Number of blocks in the longest block chain."),
			"connection_total": newGlobalMetric(namespace, "connection_total", "Number of active connections to other peers."),
			"pow_difficulty":   newGlobalMetric(namespace, "pow_difficulty", "Proof-of-work difficulty as a multiple of the minimum difficulty"),
		},
	}
}

// Describe sends the super-set of all possible descriptors of node metrics
// to the provided channel.
func (c *BitcoinCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

// Collect fetches metrics from the node and sends them to the provided channel.
func (c *BitcoinCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock() // To protect metrics from concurrent collects
	defer c.mutex.Unlock()

	stats, err := c.bitcoinClient.GetStats()
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(c.metrics["block_total"],
		prometheus.CounterValue, float64(stats.BlockCount))
	ch <- prometheus.MustNewConstMetric(c.metrics["connection_total"],
		prometheus.GaugeValue, float64(stats.ConnectionCount))
	ch <- prometheus.MustNewConstMetric(c.metrics["pow_difficulty"],
		prometheus.GaugeValue, float64(stats.Difficulty))
}
