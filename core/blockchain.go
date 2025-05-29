
package core

import (
	"errors"
	"fmt"
	"math/big"
	"sync"

	"blockchain-node/storage"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrBlockNotFound = errors.New("block not found")
	ErrInvalidBlock  = errors.New("invalid block")
)

// Blockchain represents the blockchain
type Blockchain struct {
	db           storage.Database
	currentBlock *Block
	genesis      *Block
	mu           sync.RWMutex
}

// NewBlockchain creates a new blockchain
func NewBlockchain(db storage.Database, genesis *Genesis) (*Blockchain, error) {
	bc := &Blockchain{
		db: db,
	}

	// Try to load existing blockchain
	if currentBlock, err := bc.loadCurrentBlock(); err == nil {
		bc.currentBlock = currentBlock
		if genesisBlock, err := bc.GetBlockByNumber(big.NewInt(0)); err == nil {
			bc.genesis = genesisBlock
		}
	} else {
		// Create genesis block
		genesisBlock := NewGenesisBlock(genesis)
		if err := bc.addBlock(genesisBlock); err != nil {
			return nil, fmt.Errorf("failed to add genesis block: %v", err)
		}
		bc.genesis = genesisBlock
		bc.currentBlock = genesisBlock
	}

	return bc, nil
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(block *Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Validate block
	if err := bc.validateBlock(block); err != nil {
		return fmt.Errorf("block validation failed: %v", err)
	}

	// Add to database
	if err := bc.addBlock(block); err != nil {
		return fmt.Errorf("failed to add block to database: %v", err)
	}

	bc.currentBlock = block
	return nil
}

// GetCurrentBlock returns the current (latest) block
func (bc *Blockchain) GetCurrentBlock() *Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.currentBlock
}

// GetBlockByHash retrieves a block by its hash
func (bc *Blockchain) GetBlockByHash(hash common.Hash) (*Block, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	data, err := bc.db.Get(append([]byte("block-"), hash.Bytes()...))
	if err != nil {
		return nil, ErrBlockNotFound
	}

	return deserializeBlock(data)
}

// GetBlockByNumber retrieves a block by its number
func (bc *Blockchain) GetBlockByNumber(number *big.Int) (*Block, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	// First get the hash from number index
	hashData, err := bc.db.Get(append([]byte("block-number-"), number.Bytes()...))
	if err != nil {
		return nil, ErrBlockNotFound
	}

	hash := common.BytesToHash(hashData)
	return bc.GetBlockByHash(hash)
}

// GetBlockNumber returns the current block number
func (bc *Blockchain) GetBlockNumber() *big.Int {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	
	if bc.currentBlock == nil {
		return big.NewInt(0)
	}
	return bc.currentBlock.Header.Number
}

// validateBlock validates a block
func (bc *Blockchain) validateBlock(block *Block) error {
	// Basic validation
	if block.Header.Number.Cmp(big.NewInt(0)) <= 0 && bc.currentBlock != nil {
		return ErrInvalidBlock
	}

	// Check if previous hash matches current block hash
	if bc.currentBlock != nil {
		expectedPrevHash := bc.currentBlock.Hash
		if block.Header.PreviousHash != expectedPrevHash {
			return fmt.Errorf("invalid previous hash: expected %x, got %x", 
				expectedPrevHash, block.Header.PreviousHash)
		}

		// Check block number sequence
		expectedNumber := new(big.Int).Add(bc.currentBlock.Header.Number, big.NewInt(1))
		if block.Header.Number.Cmp(expectedNumber) != 0 {
			return fmt.Errorf("invalid block number: expected %s, got %s", 
				expectedNumber.String(), block.Header.Number.String())
		}
	}

	// Validate block hash
	calculatedHash := block.CalculateHash()
	if calculatedHash != block.Hash {
		return fmt.Errorf("invalid block hash: expected %x, got %x", 
			calculatedHash, block.Hash)
	}

	return nil
}

// addBlock adds a block to the database
func (bc *Blockchain) addBlock(block *Block) error {
	// Serialize and store block
	data, err := serializeBlock(block)
	if err != nil {
		return err
	}

	// Store block by hash
	if err := bc.db.Put(append([]byte("block-"), block.Hash.Bytes()...), data); err != nil {
		return err
	}

	// Store block number index
	if err := bc.db.Put(append([]byte("block-number-"), block.Header.Number.Bytes()...), 
		block.Hash.Bytes()); err != nil {
		return err
	}

	// Update current block pointer
	if err := bc.db.Put([]byte("current-block"), block.Hash.Bytes()); err != nil {
		return err
	}

	return nil
}

// loadCurrentBlock loads the current block from database
func (bc *Blockchain) loadCurrentBlock() (*Block, error) {
	hashData, err := bc.db.Get([]byte("current-block"))
	if err != nil {
		return nil, err
	}

	hash := common.BytesToHash(hashData)
	return bc.GetBlockByHash(hash)
}

// serializeBlock serializes a block (placeholder implementation)
func serializeBlock(block *Block) ([]byte, error) {
	// TODO: Implement proper serialization (JSON/RLP)
	// For now, this is a placeholder
	return []byte(fmt.Sprintf("%+v", block)), nil
}

// deserializeBlock deserializes a block (placeholder implementation)
func deserializeBlock(data []byte) (*Block, error) {
	// TODO: Implement proper deserialization
	// For now, this is a placeholder
	return &Block{}, nil
}
