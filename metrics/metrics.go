
package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"blockchain-node/config"
	"blockchain-node/logger"
)

// Metrics holds various blockchain metrics
type Metrics struct {
	mu                    sync.RWMutex
	BlockHeight           uint64            `json:"block_height"`
	PeerCount             int               `json:"peer_count"`
	MempoolSize           int               `json:"mempool_size"`
	TotalTransactions     uint64            `json:"total_transactions"`
	MiningHashRate        float64           `json:"mining_hash_rate"`
	NetworkLatency        map[string]int64  `json:"network_latency"`
	SystemInfo            SystemInfo        `json:"system_info"`
	StartTime             time.Time         `json:"start_time"`
	logger                *logger.Logger
}

// SystemInfo holds system resource information
type SystemInfo struct {
	MemoryUsage    uint64 `json:"memory_usage_mb"`
	CPUCount       int    `json:"cpu_count"`
	GoroutineCount int    `json:"goroutine_count"`
	Uptime         string `json:"uptime"`
}

var globalMetrics *Metrics

// Init initializes the metrics system
func Init(config *config.MetricsConfig) *Metrics {
	globalMetrics = &Metrics{
		NetworkLatency: make(map[string]int64),
		StartTime:      time.Now(),
		logger:         logger.NewLogger("metrics"),
	}

	// Start metrics collection goroutine
	go globalMetrics.collectSystemMetrics()

	// Start HTTP server if enabled
	if config.Enabled {
		go globalMetrics.startHTTPServer(config)
	}

	return globalMetrics
}

// GetMetrics returns the global metrics instance
func GetMetrics() *Metrics {
	return globalMetrics
}

// UpdateBlockHeight updates the current block height
func (m *Metrics) UpdateBlockHeight(height uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.BlockHeight = height
}

// UpdatePeerCount updates the peer count
func (m *Metrics) UpdatePeerCount(count int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.PeerCount = count
}

// UpdateMempoolSize updates the mempool size
func (m *Metrics) UpdateMempoolSize(size int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MempoolSize = size
}

// IncrementTransactions increments the total transaction count
func (m *Metrics) IncrementTransactions() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalTransactions++
}

// UpdateMiningHashRate updates the mining hash rate
func (m *Metrics) UpdateMiningHashRate(hashRate float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MiningHashRate = hashRate
}

// RecordNetworkLatency records network latency to a peer
func (m *Metrics) RecordNetworkLatency(peer string, latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.NetworkLatency[peer] = latency.Milliseconds()
}

// collectSystemMetrics collects system metrics periodically
func (m *Metrics) collectSystemMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		m.mu.Lock()
		m.SystemInfo = SystemInfo{
			MemoryUsage:    memStats.Alloc / 1024 / 1024, // Convert to MB
			CPUCount:       runtime.NumCPU(),
			GoroutineCount: runtime.NumGoroutine(),
			Uptime:         time.Since(m.StartTime).String(),
		}
		m.mu.Unlock()
	}
}

// startHTTPServer starts the metrics HTTP server
func (m *Metrics) startHTTPServer(config *config.MetricsConfig) {
	mux := http.NewServeMux()
	mux.HandleFunc(config.Path, m.handleMetrics)
	
	addr := fmt.Sprintf(":%d", config.Port)
	m.logger.Info("Starting metrics server on %s%s", addr, config.Path)
	
	if err := http.ListenAndServe(addr, mux); err != nil {
		m.logger.Error("Metrics server error: %v", err)
	}
}

// handleMetrics handles HTTP requests for metrics
func (m *Metrics) handleMetrics(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(m); err != nil {
		m.logger.Error("Failed to encode metrics: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// GetSnapshot returns a copy of current metrics
func (m *Metrics) GetSnapshot() Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a deep copy
	snapshot := *m
	snapshot.NetworkLatency = make(map[string]int64)
	for k, v := range m.NetworkLatency {
		snapshot.NetworkLatency[k] = v
	}

	return snapshot
}
