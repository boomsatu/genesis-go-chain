
package consensus

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"blockchain-node/core"

	"github.com/ethereum/go-ethereum/common"
)

// ProofOfWork represents the Proof of Work consensus engine
type ProofOfWork struct {
	difficulty *big.Int
}

// NewProofOfWork creates a new PoW instance
func NewProofOfWork(difficulty *big.Int) *ProofOfWork {
	return &ProofOfWork{
		difficulty: difficulty,
	}
}

// Mine mines a block using Proof of Work
func (pow *ProofOfWork) Mine(block *core.Block) error {
	fmt.Printf("Mining block with difficulty %s...\n", pow.difficulty.String())
	
	start := time.Now()
	nonce := uint64(0)
	target := pow.calculateTarget()
	
	for {
		// Update nonce in block header
		block.Header.Nonce = nonce
		
		// Calculate hash
		hash := pow.calculateHash(block)
		hashInt := new(big.Int).SetBytes(hash[:])
		
		// Check if hash meets difficulty target
		if hashInt.Cmp(target) == -1 {
			// Found valid hash
			block.Hash = hash
			elapsed := time.Since(start)
			fmt.Printf("Block mined! Nonce: %d, Hash: %x, Time: %v\n", 
				nonce, hash, elapsed)
			return nil
		}
		
		nonce++
		
		// Progress indicator
		if nonce%100000 == 0 {
			fmt.Printf("Mining... nonce: %d\n", nonce)
		}
	}
}

// ValidateBlock validates a block's proof of work
func (pow *ProofOfWork) ValidateBlock(block *core.Block) bool {
	target := pow.calculateTarget()
	hash := pow.calculateHash(block)
	hashInt := new(big.Int).SetBytes(hash[:])
	
	return hashInt.Cmp(target) == -1
}

// calculateTarget calculates the target value for mining
func (pow *ProofOfWork) calculateTarget() *big.Int {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-pow.difficulty.Uint64()))
	return target
}

// calculateHash calculates the hash for a block
func (pow *ProofOfWork) calculateHash(block *core.Block) common.Hash {
	// Combine header data for hashing
	data := append(block.Header.PreviousHash.Bytes(), block.Header.StateRoot.Bytes()...)
	data = append(data, block.Header.TransactionsRoot.Bytes()...)
	data = append(data, block.Header.Number.Bytes()...)
	data = append(data, big.NewInt(int64(block.Header.Timestamp)).Bytes()...)
	data = append(data, big.NewInt(int64(block.Header.Nonce)).Bytes()...)
	data = append(data, pow.difficulty.Bytes()...)
	
	hash := sha256.Sum256(data)
	return common.BytesToHash(hash[:])
}

// SetDifficulty updates the mining difficulty
func (pow *ProofOfWork) SetDifficulty(difficulty *big.Int) {
	pow.difficulty = difficulty
}

// GetDifficulty returns the current difficulty
func (pow *ProofOfWork) GetDifficulty() *big.Int {
	return new(big.Int).Set(pow.difficulty)
}
