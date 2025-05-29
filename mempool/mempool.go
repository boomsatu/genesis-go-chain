
package mempool

import (
	"errors"
	"fmt"
	"sync"

	"blockchain-node/core"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ErrTransactionExists     = errors.New("transaction already exists")
	ErrInvalidTransaction    = errors.New("invalid transaction")
	ErrInsufficientGasPrice = errors.New("gas price too low")
)

// Mempool represents the transaction pool
type Mempool struct {
	transactions map[common.Hash]*core.Transaction
	pending      map[common.Address][]*core.Transaction // Pending transactions by sender
	mu           sync.RWMutex
	minGasPrice  uint64
}

// NewMempool creates a new mempool
func NewMempool() *Mempool {
	return &Mempool{
		transactions: make(map[common.Hash]*core.Transaction),
		pending:      make(map[common.Address][]*core.Transaction),
		minGasPrice:  1000000000, // 1 Gwei
	}
}

// AddTransaction adds a transaction to the mempool
func (mp *Mempool) AddTransaction(tx *core.Transaction) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	// Check if transaction already exists
	if _, exists := mp.transactions[tx.Hash]; exists {
		return ErrTransactionExists
	}

	// Validate transaction
	if err := mp.validateTransaction(tx); err != nil {
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	// Add to mempool
	mp.transactions[tx.Hash] = tx
	mp.pending[tx.From] = append(mp.pending[tx.From], tx)

	fmt.Printf("Transaction added to mempool: %x\n", tx.Hash)
	return nil
}

// RemoveTransaction removes a transaction from the mempool
func (mp *Mempool) RemoveTransaction(txHash common.Hash) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if tx, exists := mp.transactions[txHash]; exists {
		delete(mp.transactions, txHash)
		
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

// GetPendingTransactionsForMining returns transactions ready for mining
func (mp *Mempool) GetPendingTransactionsForMining(limit int) []*core.Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	var txs []*core.Transaction
	count := 0
	
	// Simple selection - in production, should sort by gas price and nonce
	for _, tx := range mp.transactions {
		if count >= limit {
			break
		}
		txs = append(txs, tx)
		count++
	}
	
	return txs
}

// Size returns the number of transactions in the mempool
func (mp *Mempool) Size() int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	return len(mp.transactions)
}

// validateTransaction validates a transaction
func (mp *Mempool) validateTransaction(tx *core.Transaction) error {
	// Check gas price
	if tx.GasPrice.Uint64() < mp.minGasPrice {
		return ErrInsufficientGasPrice
	}

	// Check basic transaction validity
	if tx.GasLimit == 0 {
		return fmt.Errorf("gas limit cannot be zero")
	}

	if tx.Value.Sign() < 0 {
		return fmt.Errorf("negative value")
	}

	// Validate signature
	if err := mp.validateSignature(tx); err != nil {
		return fmt.Errorf("invalid signature: %v", err)
	}

	return nil
}

// validateSignature validates the transaction signature
func (mp *Mempool) validateSignature(tx *core.Transaction) error {
	// Recover sender from signature
	hash := tx.CalculateHash()
	
	// Simple signature validation - in production, use proper ECDSA recovery
	if tx.V == nil || tx.R == nil || tx.S == nil {
		return fmt.Errorf("missing signature components")
	}

	// For now, just check that signature components exist
	// TODO: Implement proper ECDSA signature verification
	_ = hash // Use hash for signature verification
	
	return nil
}

// recoverSender recovers the sender address from transaction signature
func (mp *Mempool) recoverSender(tx *core.Transaction) (common.Address, error) {
	// Create signing hash
	hash := tx.CalculateHash()
	
	// TODO: Implement proper ECDSA recovery
	// This is a placeholder that would need proper implementation
	_ = hash
	
	return tx.From, nil
}

// SetMinGasPrice sets the minimum gas price for transactions
func (mp *Mempool) SetMinGasPrice(gasPrice uint64) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.minGasPrice = gasPrice
}

// GetMinGasPrice returns the minimum gas price
func (mp *Mempool) GetMinGasPrice() uint64 {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	return mp.minGasPrice
}

// Clear removes all transactions from the mempool
func (mp *Mempool) Clear() {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	
	mp.transactions = make(map[common.Hash]*core.Transaction)
	mp.pending = make(map[common.Address][]*core.Transaction)
}
