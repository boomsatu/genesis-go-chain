
package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"blockchain-node/config"
	"blockchain-node/logger"

	"github.com/gorilla/mux"
)

// Metrics holds all blockchain metrics
type Metrics struct {
	config  *config.MetricsConfig
	logger  *logger.Logger
	server  *http.Server
	mu      sync.RWMutex
	
	// Blockchain metrics
	BlockHeight       uint64    `json:"block_height"`
	TotalTransactions uint64    `json:"total_transactions"`
	MempoolSize       int       `json:"mempool_size"`
	PeerCount         int       `json:"peer_count"`
	
	// Mining metrics
	HashRate          float64   `json:"hash_rate"`
	BlocksMinedCount  uint64    `json:"blocks_mined_count"`
	MiningDifficulty  uint64    `json:"mining_difficulty"`
	
	// Performance metrics
	BlockProcessingTime time.Duration `json:"block_processing_time_ns"`
	TxProcessingTime    time.Duration `json:"tx_processing_time_ns"`
	DatabaseSize        uint64        `json:"database_size_bytes"`
	
	// Network metrics
	InboundConnections  int `json:"inbound_connections"`
	OutboundConnections int `json:"outbound_connections"`
	MessagesSent        uint64 `json:"messages_sent"`
	MessagesReceived    uint64 `json:"messages_received"`
	
	// System metrics
	StartTime         time.Time `json:"start_time"`
	Uptime            time.Duration `json:"uptime_seconds"`
	MemoryUsage       uint64    `json:"memory_usage_bytes"`
	CPUUsage          float64   `json:"cpu_usage_percent"`
	
	// Custom metrics
	CustomMetrics map[string]interface{} `json:"custom_metrics"`
}

// Init initializes the metrics system
func Init(config *config.MetricsConfig) *Metrics {
	metrics := &Metrics{
		config:        config,
		logger:        logger.NewLogger("metrics"),
		StartTime:     time.Now(),
		CustomMetrics: make(map[string]interface{}),
	}

	if config.Enabled {
		if err := metrics.startServer(); err != nil {
			metrics.logger.Error("Failed to start metrics server", "error", err)
		}
	}

	metrics.logger.Info("Metrics system initialized", "enabled", config.Enabled)
	return metrics
}

// startServer starts the metrics HTTP server
func (m *Metrics) startServer() error {
	router := mux.NewRouter()
	
	// Metrics endpoint
	router.HandleFunc(m.config.Path, m.handleMetrics).Methods("GET")
	
	// Prometheus-style metrics endpoint
	router.HandleFunc("/metrics", m.handlePrometheusMetrics).Methods("GET")
	
	// Health endpoint
	router.HandleFunc("/health", m.handleHealth).Methods("GET")

	m.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", m.config.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		m.logger.Info("Starting metrics server", "port", m.config.Port, "path", m.config.Path)
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			m.logger.Error("Metrics server error", "error", err)
		}
	}()

	return nil
}

// Stop stops the metrics server
func (m *Metrics) Stop() error {
	if m.server != nil {
		m.logger.Info("Stopping metrics server...")
		return m.server.Close()
	}
	return nil
}

// handleMetrics handles the JSON metrics endpoint
func (m *Metrics) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	m.mu.RLock()
	// Update uptime
	m.Uptime = time.Since(m.StartTime)
	
	// Create a copy for safe JSON marshaling
	metricsCopy := *m
	m.mu.RUnlock()

	if err := json.NewEncoder(w).Encode(metricsCopy); err != nil {
		m.logger.Error("Failed to encode metrics", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// handlePrometheusMetrics handles Prometheus-style metrics
func (m *Metrics) handlePrometheusMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Update uptime
	uptime := time.Since(m.StartTime).Seconds()

	fmt.Fprintf(w, "# HELP lumina_block_height Current block height\n")
	fmt.Fprintf(w, "# TYPE lumina_block_height gauge\n")
	fmt.Fprintf(w, "lumina_block_height %d\n", m.BlockHeight)

	fmt.Fprintf(w, "# HELP lumina_total_transactions Total number of transactions processed\n")
	fmt.Fprintf(w, "# TYPE lumina_total_transactions counter\n")
	fmt.Fprintf(w, "lumina_total_transactions %d\n", m.TotalTransactions)

	fmt.Fprintf(w, "# HELP lumina_mempool_size Current mempool size\n")
	fmt.Fprintf(w, "# TYPE lumina_mempool_size gauge\n")
	fmt.Fprintf(w, "lumina_mempool_size %d\n", m.MempoolSize)

	fmt.Fprintf(w, "# HELP lumina_peer_count Number of connected peers\n")
	fmt.Fprintf(w, "# TYPE lumina_peer_count gauge\n")
	fmt.Fprintf(w, "lumina_peer_count %d\n", m.PeerCount)

	fmt.Fprintf(w, "# HELP lumina_hash_rate Current mining hash rate\n")
	fmt.Fprintf(w, "# TYPE lumina_hash_rate gauge\n")
	fmt.Fprintf(w, "lumina_hash_rate %f\n", m.HashRate)

	fmt.Fprintf(w, "# HELP lumina_blocks_mined_total Total blocks mined\n")
	fmt.Fprintf(w, "# TYPE lumina_blocks_mined_total counter\n")
	fmt.Fprintf(w, "lumina_blocks_mined_total %d\n", m.BlocksMinedCount)

	fmt.Fprintf(w, "# HELP lumina_uptime_seconds Node uptime in seconds\n")
	fmt.Fprintf(w, "# TYPE lumina_uptime_seconds gauge\n")
	fmt.Fprintf(w, "lumina_uptime_seconds %f\n", uptime)

	fmt.Fprintf(w, "# HELP lumina_block_processing_time_seconds Time to process last block\n")
	fmt.Fprintf(w, "# TYPE lumina_block_processing_time_seconds gauge\n")
	fmt.Fprintf(w, "lumina_block_processing_time_seconds %f\n", m.BlockProcessingTime.Seconds())

	fmt.Fprintf(w, "# HELP lumina_messages_sent_total Total messages sent to peers\n")
	fmt.Fprintf(w, "# TYPE lumina_messages_sent_total counter\n")
	fmt.Fprintf(w, "lumina_messages_sent_total %d\n", m.MessagesSent)

	fmt.Fprintf(w, "# HELP lumina_messages_received_total Total messages received from peers\n")
	fmt.Fprintf(w, "# TYPE lumina_messages_received_total counter\n")
	fmt.Fprintf(w, "lumina_messages_received_total %d\n", m.MessagesReceived)
}

// handleHealth handles health check requests
func (m *Metrics) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"uptime":    time.Since(m.StartTime).Seconds(),
	}

	json.NewEncoder(w).Encode(health)
}

// Update methods for various metrics

func (m *Metrics) UpdateBlockHeight(height uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.BlockHeight = height
}

func (m *Metrics) IncrementTransactions() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalTransactions++
}

func (m *Metrics) UpdateMempoolSize(size int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MempoolSize = size
}

func (m *Metrics) UpdatePeerCount(count int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.PeerCount = count
}

func (m *Metrics) UpdateMiningHashRate(hashRate float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.HashRate = hashRate
}

func (m *Metrics) IncrementBlocksMined() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.BlocksMinedCount++
}

func (m *Metrics) UpdateMiningDifficulty(difficulty uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MiningDifficulty = difficulty
}

func (m *Metrics) UpdateBlockProcessingTime(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.BlockProcessingTime = duration
}

func (m *Metrics) UpdateTxProcessingTime(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TxProcessingTime = duration
}

func (m *Metrics) UpdateDatabaseSize(size uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.DatabaseSize = size
}

func (m *Metrics) UpdateNetworkConnections(inbound, outbound int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.InboundConnections = inbound
	m.OutboundConnections = outbound
}

func (m *Metrics) IncrementMessagesSent() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MessagesSent++
}

func (m *Metrics) IncrementMessagesReceived() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MessagesReceived++
}

func (m *Metrics) UpdateMemoryUsage(usage uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MemoryUsage = usage
}

func (m *Metrics) UpdateCPUUsage(usage float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CPUUsage = usage
}

func (m *Metrics) SetCustomMetric(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CustomMetrics[key] = value
}

func (m *Metrics) GetCustomMetric(key string) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.CustomMetrics[key]
}

// GetSnapshot returns a snapshot of current metrics
func (m *Metrics) GetSnapshot() *Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a deep copy
	snapshot := *m
	snapshot.Uptime = time.Since(m.StartTime)
	
	// Copy custom metrics map
	snapshot.CustomMetrics = make(map[string]interface{})
	for k, v := range m.CustomMetrics {
		snapshot.CustomMetrics[k] = v
	}

	return &snapshot
}

// Reset resets all metrics to zero (useful for testing)
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.BlockHeight = 0
	m.TotalTransactions = 0
	m.MempoolSize = 0
	m.PeerCount = 0
	m.HashRate = 0
	m.BlocksMinedCount = 0
	m.MiningDifficulty = 0
	m.BlockProcessingTime = 0
	m.TxProcessingTime = 0
	m.DatabaseSize = 0
	m.InboundConnections = 0
	m.OutboundConnections = 0
	m.MessagesSent = 0
	m.MessagesReceived = 0
	m.StartTime = time.Now()
	m.MemoryUsage = 0
	m.CPUUsage = 0
	m.CustomMetrics = make(map[string]interface{})

	m.logger.Info("Metrics reset")
}

// LogMetrics logs current metrics at INFO level
func (m *Metrics) LogMetrics() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.logger.Info("Current metrics snapshot",
		"block_height", m.BlockHeight,
		"total_transactions", m.TotalTransactions,
		"mempool_size", m.MempoolSize,
		"peer_count", m.PeerCount,
		"hash_rate", m.HashRate,
		"blocks_mined", m.BlocksMinedCount,
		"uptime_seconds", time.Since(m.StartTime).Seconds(),
	)
}
