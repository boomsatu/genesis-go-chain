
package mempool

import (
	"container/heap"
	"fmt"
	"math/big"
	"sync"
	"time"

	"blockchain-node/core"
	"blockchain-node/logger"

	"github.com/ethereum/go-ethereum/common"
)

// Config holds mempool configuration
type Config struct {
	MaxSize     int      // Maximum number of transactions
	MinGasPrice uint64   // Minimum gas price (wei)
	MaxTxSize   int      // Maximum transaction size in bytes
	Timeout     duration // Transaction timeout
}

type duration time.Duration

// Mempool manages pending transactions
type Mempool struct {
	config      *Config
	pending     map[common.Hash]*core.Transaction
	queue       TransactionQueue
	byFrom      map[common.Address][]*core.Transaction
	logger      *logger.Logger
	mu          sync.RWMutex
}

// TransactionPriorityItem represents a transaction with priority for the heap
type TransactionPriorityItem struct {
	Tx       *core.Transaction
	Priority *big.Int // Gas price for priority
	Index    int
}

// TransactionQueue implements heap.Interface for transaction prioritization
type TransactionQueue []*TransactionPriorityItem

func (pq TransactionQueue) Len() int { return len(pq) }

func (pq TransactionQueue) Less(i, j int) bool {
	// Higher gas price has higher priority
	return pq[i].Priority.Cmp(pq[j].Priority) > 0
}

func (pq TransactionQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *TransactionQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*TransactionPriorityItem)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *TransactionQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

// NewMempool creates a new mempool instance
func NewMempool(config *Config) *Mempool {
	return &Mempool{
		config:  config,
		pending: make(map[common.Hash]*core.Transaction),
		queue:   make(TransactionQueue, 0),
		byFrom:  make(map[common.Address][]*core.Transaction),
		logger:  logger.NewLogger("mempool"),
	}
}

// AddTransaction adds a transaction to the mempool
func (mp *Mempool) AddTransaction(tx *core.Transaction) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	// Validate transaction
	if err := mp.validateTransaction(tx); err != nil {
		mp.logger.Warning("Transaction validation failed", "hash", tx.Hash.Hex(), "error", err)
		return err
	}

	// Check if transaction already exists
	if _, exists := mp.pending[tx.Hash]; exists {
		return fmt.Errorf("transaction already exists in mempool")
	}

	// Check mempool size limit
	if len(mp.pending) >= mp.config.MaxSize {
		// Remove lowest priority transaction
		mp.removeLowPriorityTransaction()
	}

	// Add to pending transactions
	mp.pending[tx.Hash] = tx

	// Add to priority queue
	item := &TransactionPriorityItem{
		Tx:       tx,
		Priority: tx.GasPrice,
	}
	heap.Push(&mp.queue, item)

	// Add to by-from index
	mp.byFrom[tx.From] = append(mp.byFrom[tx.From], tx)

	mp.logger.Debug("Transaction added to mempool", 
		"hash", tx.Hash.Hex(), 
		"from", tx.From.Hex(), 
		"gasPrice", tx.GasPrice.String(),
		"mempoolSize", len(mp.pending))

	return nil
}

// RemoveTransaction removes a transaction from the mempool
func (mp *Mempool) RemoveTransaction(hash common.Hash) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	tx, exists := mp.pending[hash]
	if !exists {
		return
	}

	// Remove from pending
	delete(mp.pending, hash)

	// Remove from by-from index
	fromTxs := mp.byFrom[tx.From]
	for i, fromTx := range fromTxs {
		if fromTx.Hash == hash {
			mp.byFrom[tx.From] = append(fromTxs[:i], fromTxs[i+1:]...)
			break
		}
	}

	// Remove empty slice
	if len(mp.byFrom[tx.From]) == 0 {
		delete(mp.byFrom, tx.From)
	}

	// Rebuild priority queue (inefficient but simple)
	mp.rebuildQueue()

	mp.logger.Debug("Transaction removed from mempool", 
		"hash", hash.Hex(), 
		"mempoolSize", len(mp.pending))
}

// GetTransaction retrieves a transaction by hash
func (mp *Mempool) GetTransaction(hash common.Hash) *core.Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return mp.pending[hash]
}

// GetPendingTransactions returns all pending transactions
func (mp *Mempool) GetPendingTransactions() []*core.Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	txs := make([]*core.Transaction, 0, len(mp.pending))
	for _, tx := range mp.pending {
		txs = append(txs, tx)
	}

	return txs
}

// GetPendingTransactionsForMining returns transactions ready for mining
func (mp *Mempool) GetPendingTransactionsForMining(maxCount int) []*core.Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	if len(mp.queue) == 0 {
		return []*core.Transaction{}
	}

	// Create a copy of the queue for processing
	queueCopy := make(TransactionQueue, len(mp.queue))
	copy(queueCopy, mp.queue)
	heap.Init(&queueCopy)

	txs := make([]*core.Transaction, 0, maxCount)
	count := 0

	for len(queueCopy) > 0 && count < maxCount {
		item := heap.Pop(&queueCopy).(*TransactionPriorityItem)
		txs = append(txs, item.Tx)
		count++
	}

	return txs
}

// GetTransactionsByFrom returns transactions from a specific address
func (mp *Mempool) GetTransactionsByFrom(from common.Address) []*core.Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	txs := mp.byFrom[from]
	if txs == nil {
		return []*core.Transaction{}
	}

	// Return a copy
	result := make([]*core.Transaction, len(txs))
	copy(result, txs)
	return result
}

// Size returns the current size of the mempool
func (mp *Mempool) Size() int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return len(mp.pending)
}

// validateTransaction validates a transaction before adding to mempool
func (mp *Mempool) validateTransaction(tx *core.Transaction) error {
	// Check minimum gas price
	if tx.GasPrice.Cmp(big.NewInt(int64(mp.config.MinGasPrice))) < 0 {
		return fmt.Errorf("gas price too low: got %s, minimum %d", 
			tx.GasPrice.String(), mp.config.MinGasPrice)
	}

	// Check gas limit
	if tx.GasLimit == 0 {
		return fmt.Errorf("gas limit cannot be zero")
	}

	if tx.GasLimit > 8000000 { // Max block gas limit
		return fmt.Errorf("gas limit too high: %d", tx.GasLimit)
	}

	// Check transaction size
	if mp.config.MaxTxSize > 0 {
		// Estimate transaction size (simplified)
		txSize := 32 + 32 + 8 + 8 + 32 + len(tx.Data) + 32 + 32 + 32 // Basic fields + data + signature
		if txSize > mp.config.MaxTxSize {
			return fmt.Errorf("transaction too large: %d bytes", txSize)
		}
	}

	// Check for valid signature components
	if tx.V == nil || tx.R == nil || tx.S == nil {
		return fmt.Errorf("invalid signature components")
	}

	// Basic value validation
	if tx.Value == nil {
		return fmt.Errorf("value cannot be nil")
	}

	if tx.Value.Sign() < 0 {
		return fmt.Errorf("negative value not allowed")
	}

	return nil
}

// removeLowPriorityTransaction removes the transaction with lowest priority
func (mp *Mempool) removeLowPriorityTransaction() {
	if len(mp.queue) == 0 {
		return
	}

	// Find transaction with lowest gas price
	var lowestTx *core.Transaction
	lowestGasPrice := new(big.Int)

	for _, tx := range mp.pending {
		if lowestTx == nil || tx.GasPrice.Cmp(lowestGasPrice) < 0 {
			lowestTx = tx
			lowestGasPrice = tx.GasPrice
		}
	}

	if lowestTx != nil {
		mp.logger.Debug("Removing low priority transaction", 
			"hash", lowestTx.Hash.Hex(), 
			"gasPrice", lowestTx.GasPrice.String())
		
		// Remove without locking (already locked)
		delete(mp.pending, lowestTx.Hash)
		
		// Remove from by-from index
		fromTxs := mp.byFrom[lowestTx.From]
		for i, fromTx := range fromTxs {
			if fromTx.Hash == lowestTx.Hash {
				mp.byFrom[lowestTx.From] = append(fromTxs[:i], fromTxs[i+1:]...)
				break
			}
		}

		if len(mp.byFrom[lowestTx.From]) == 0 {
			delete(mp.byFrom, lowestTx.From)
		}

		mp.rebuildQueue()
	}
}

// rebuildQueue rebuilds the priority queue
func (mp *Mempool) rebuildQueue() {
	mp.queue = make(TransactionQueue, 0, len(mp.pending))
	
	for _, tx := range mp.pending {
		item := &TransactionPriorityItem{
			Tx:       tx,
			Priority: tx.GasPrice,
		}
		mp.queue = append(mp.queue, item)
	}

	heap.Init(&mp.queue)
}

// Clean removes expired transactions from mempool
func (mp *Mempool) Clean() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	// For now, we don't implement timeout-based cleaning
	// This could be added based on transaction timestamp vs current time
	mp.logger.Debug("Mempool cleanup completed", "size", len(mp.pending))
}

// GetStats returns mempool statistics
func (mp *Mempool) GetStats() map[string]interface{} {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	stats := map[string]interface{}{
		"pending_count":  len(mp.pending),
		"queue_length":   len(mp.queue),
		"unique_senders": len(mp.byFrom),
		"max_size":       mp.config.MaxSize,
		"min_gas_price":  mp.config.MinGasPrice,
	}

	// Calculate average gas price
	if len(mp.pending) > 0 {
		totalGasPrice := big.NewInt(0)
		for _, tx := range mp.pending {
			totalGasPrice.Add(totalGasPrice, tx.GasPrice)
		}
		avgGasPrice := new(big.Int).Div(totalGasPrice, big.NewInt(int64(len(mp.pending))))
		stats["avg_gas_price"] = avgGasPrice.String()
	}

	return stats
}

// GetTransactionHashes returns all transaction hashes in mempool
func (mp *Mempool) GetTransactionHashes() []common.Hash {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	hashes := make([]common.Hash, 0, len(mp.pending))
	for hash := range mp.pending {
		hashes = append(hashes, hash)
	}

	return hashes
}

// HasTransaction checks if a transaction exists in mempool
func (mp *Mempool) HasTransaction(hash common.Hash) bool {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	_, exists := mp.pending[hash]
	return exists
}

// GetHighestGasPriceTransaction returns the transaction with highest gas price
func (mp *Mempool) GetHighestGasPriceTransaction() *core.Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	if len(mp.queue) == 0 {
		return nil
	}

	return mp.queue[0].Tx
}
