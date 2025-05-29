
package mempool

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"blockchain-node/core"
	"blockchain-node/logger"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrTransactionExists     = errors.New("transaction already exists")
	ErrInvalidTransaction    = errors.New("invalid transaction")
	ErrInsufficientGasPrice = errors.New("gas price too low")
	ErrMempoolFull          = errors.New("mempool is full")
	ErrInvalidNonce         = errors.New("invalid nonce")
)

// Config holds mempool configuration
type Config struct {
	MaxSize     int    `json:"max_size"`
	MinGasPrice uint64 `json:"min_gas_price"`
}

// Mempool represents the transaction pool
type Mempool struct {
	transactions map[common.Hash]*core.Transaction
	pending      map[common.Address][]*core.Transaction // Pending transactions by sender
	mu           sync.RWMutex
	config       *Config
	logger       *logger.Logger
	
	// Statistics
	totalAdded    uint64
	totalRemoved  uint64
	totalRejected uint64
}

// NewMempool creates a new mempool with configuration
func NewMempool(config *Config) *Mempool {
	if config == nil {
		config = &Config{
			MaxSize:     1000,
			MinGasPrice: 1000000000, // 1 Gwei
		}
	}

	return &Mempool{
		transactions: make(map[common.Hash]*core.Transaction),
		pending:      make(map[common.Address][]*core.Transaction),
		config:       config,
		logger:       logger.NewLogger("mempool"),
	}
}

// AddTransaction adds a transaction to the mempool
func (mp *Mempool) AddTransaction(tx *core.Transaction) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	// Check if mempool is full
	if len(mp.transactions) >= mp.config.MaxSize {
		mp.totalRejected++
		return ErrMempoolFull
	}

	// Check if transaction already exists
	if _, exists := mp.transactions[tx.Hash]; exists {
		mp.totalRejected++
		return ErrTransactionExists
	}

	// Validate transaction
	if err := mp.validateTransaction(tx); err != nil {
		mp.totalRejected++
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	// Add to mempool
	mp.transactions[tx.Hash] = tx
	mp.pending[tx.From] = append(mp.pending[tx.From], tx)
	mp.totalAdded++

	mp.logger.Debug("Transaction added to mempool: %x, from: %s, gas: %d", 
		tx.Hash, tx.From.Hex(), tx.GasLimit)
	
	return nil
}

// RemoveTransaction removes a transaction from the mempool
func (mp *Mempool) RemoveTransaction(txHash common.Hash) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if tx, exists := mp.transactions[txHash]; exists {
		delete(mp.transactions, txHash)
		mp.totalRemoved++
		
		// Remove from pending list
		senderTxs := mp.pending[tx.From]
		for i, pendingTx := range senderTxs {
			if pendingTx.Hash == txHash {
				mp.pending[tx.From] = append(senderTxs[:i], senderTxs[i+1:]...)
				break
			}
		}
		
		// Clean up empty pending lists
		if len(mp.pending[tx.From]) == 0 {
			delete(mp.pending, tx.From)
		}

		mp.logger.Debug("Transaction removed from mempool: %x", txHash)
	}
}

// GetTransaction retrieves a transaction by hash
func (mp *Mempool) GetTransaction(txHash common.Hash) (*core.Transaction, bool) {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	tx, exists := mp.transactions[txHash]
	return tx, exists
}

// GetPendingTransactions returns all pending transactions
func (mp *Mempool) GetPendingTransactions() []*core.Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	var txs []*core.Transaction
	for _, tx := range mp.transactions {
		txs = append(txs, tx)
	}
	return txs
}

// GetPendingTransactionsForMining returns transactions ready for mining, sorted by gas price
func (mp *Mempool) GetPendingTransactionsForMining(limit int) []*core.Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	// Convert to slice
	var txs []*core.Transaction
	for _, tx := range mp.transactions {
		txs = append(txs, tx)
	}

	// Sort by gas price (descending) - simple bubble sort for now
	// In production, use a more efficient sorting algorithm
	for i := 0; i < len(txs)-1; i++ {
		for j := 0; j < len(txs)-i-1; j++ {
			if txs[j].GasPrice.Cmp(txs[j+1].GasPrice) < 0 {
				txs[j], txs[j+1] = txs[j+1], txs[j]
			}
		}
	}

	// Apply limit
	if limit > 0 && len(txs) > limit {
		txs = txs[:limit]
	}

	mp.logger.Debug("Retrieved %d transactions for mining (limit: %d)", len(txs), limit)
	return txs
}

// Size returns the number of transactions in the mempool
func (mp *Mempool) Size() int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	return len(mp.transactions)
}

// GetPendingTransactionsByAddress returns pending transactions for a specific address
func (mp *Mempool) GetPendingTransactionsByAddress(address common.Address) []*core.Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	txs := mp.pending[address]
	result := make([]*core.Transaction, len(txs))
	copy(result, txs)
	return result
}

// validateTransaction validates a transaction
func (mp *Mempool) validateTransaction(tx *core.Transaction) error {
	// Check gas price
	if tx.GasPrice.Uint64() < mp.config.MinGasPrice {
		return fmt.Errorf("%w: got %d, minimum %d", ErrInsufficientGasPrice, 
			tx.GasPrice.Uint64(), mp.config.MinGasPrice)
	}

	// Check basic transaction validity
	if tx.GasLimit == 0 {
		return fmt.Errorf("gas limit cannot be zero")
	}

	if tx.Value.Sign() < 0 {
		return fmt.Errorf("negative value")
	}

	// Check transaction size limit (prevent spam)
	if len(tx.Data) > 32*1024 { // 32KB limit
		return fmt.Errorf("transaction data too large: %d bytes", len(tx.Data))
	}

	// Validate signature
	if err := mp.validateSignature(tx); err != nil {
		return fmt.Errorf("invalid signature: %v", err)
	}

	return nil
}

// validateSignature validates the transaction signature
func (mp *Mempool) validateSignature(tx *core.Transaction) error {
	// Simple signature validation - in production, use proper ECDSA recovery
	if tx.V == nil || tx.R == nil || tx.S == nil {
		return fmt.Errorf("missing signature components")
	}

	// TODO: Implement proper ECDSA signature verification
	// For now, just check that signature components exist and are valid
	if tx.V.Sign() < 0 || tx.R.Sign() <= 0 || tx.S.Sign() <= 0 {
		return fmt.Errorf("invalid signature values")
	}
	
	return nil
}

// SetMinGasPrice sets the minimum gas price for transactions
func (mp *Mempool) SetMinGasPrice(gasPrice uint64) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	
	oldPrice := mp.config.MinGasPrice
	mp.config.MinGasPrice = gasPrice
	
	mp.logger.Info("Minimum gas price updated: %d -> %d", oldPrice, gasPrice)
	
	// Remove transactions with gas price below new minimum
	var toRemove []common.Hash
	for hash, tx := range mp.transactions {
		if tx.GasPrice.Uint64() < gasPrice {
			toRemove = append(toRemove, hash)
		}
	}
	
	for _, hash := range toRemove {
		mp.RemoveTransaction(hash)
		mp.logger.Debug("Removed transaction due to low gas price: %x", hash)
	}
}

// GetMinGasPrice returns the minimum gas price
func (mp *Mempool) GetMinGasPrice() uint64 {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	return mp.config.MinGasPrice
}

// Clear removes all transactions from the mempool
func (mp *Mempool) Clear() {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	
	count := len(mp.transactions)
	mp.transactions = make(map[common.Hash]*core.Transaction)
	mp.pending = make(map[common.Address][]*core.Transaction)
	
	mp.logger.Info("Mempool cleared: %d transactions removed", count)
}

// GetStats returns mempool statistics
func (mp *Mempool) GetStats() map[string]interface{} {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return map[string]interface{}{
		"current_size":    len(mp.transactions),
		"max_size":        mp.config.MaxSize,
		"min_gas_price":   mp.config.MinGasPrice,
		"total_added":     mp.totalAdded,
		"total_removed":   mp.totalRemoved,
		"total_rejected":  mp.totalRejected,
		"pending_senders": len(mp.pending),
	}
}

// Cleanup removes old or invalid transactions periodically
func (mp *Mempool) Cleanup() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	// Remove transactions older than 1 hour (simple cleanup)
	// In production, implement more sophisticated cleanup logic
	cutoffTime := time.Now().Add(-1 * time.Hour).Unix()
	var toRemove []common.Hash
	
	for hash, tx := range mp.transactions {
		// Assume transaction has a timestamp field (would need to be added to Transaction struct)
		// For now, just demonstrate the cleanup pattern
		_ = cutoffTime
		_ = tx
		// if tx.Timestamp < cutoffTime {
		//     toRemove = append(toRemove, hash)
		// }
	}
	
	for _, hash := range toRemove {
		if tx := mp.transactions[hash]; tx != nil {
			delete(mp.transactions, hash)
			mp.totalRemoved++
			mp.logger.Debug("Cleaned up old transaction: %x", hash)
		}
	}
}
