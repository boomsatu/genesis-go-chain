
package core

import (
	"errors"
	"fmt"
	"math/big"

	"blockchain-node/crypto"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidNonce        = errors.New("invalid nonce")
	ErrGasLimitExceeded    = errors.New("gas limit exceeded")
	ErrInvalidSignature    = errors.New("invalid signature")
)

// ExecutionEngine represents the custom transaction execution environment
type ExecutionEngine struct {
	stateDB *StateDB
	config  *ExecutionConfig
}

// ExecutionConfig holds configuration for the execution engine
type ExecutionConfig struct {
	ChainID       *big.Int
	BlockGasLimit uint64
	MinGasPrice   *big.Int
}

// ExecutionResult contains the result of transaction execution
type ExecutionResult struct {
	GasUsed         uint64
	Status          uint64 // 1 for success, 0 for failure
	Logs            []*Log
	ContractAddress *crypto.Address // For contract creation
	Error           error
}

// NewExecutionEngine creates a new execution engine
func NewExecutionEngine(stateDB *StateDB, config *ExecutionConfig) *ExecutionEngine {
	return &ExecutionEngine{
		stateDB: stateDB,
		config:  config,
	}
}

// ExecuteTransaction executes a transaction in the custom environment
func (ee *ExecutionEngine) ExecuteTransaction(tx *Transaction, header *BlockHeader) (*ExecutionResult, error) {
	// Validate transaction signature
	if err := ee.validateSignature(tx); err != nil {
		return &ExecutionResult{Status: 0, Error: err}, err
	}

	// Get sender account
	senderAccount := ee.stateDB.GetAccount(tx.From)
	if senderAccount == nil {
		senderAccount = &Account{
			Nonce:   0,
			Balance: big.NewInt(0),
		}
	}

	// Validate nonce
	if senderAccount.Nonce != tx.Nonce {
		return &ExecutionResult{Status: 0, Error: ErrInvalidNonce}, ErrInvalidNonce
	}

	// Calculate total cost (value + gas)
	gasCost := new(big.Int).Mul(tx.GasPrice, big.NewInt(int64(tx.GasLimit)))
	totalCost := new(big.Int).Add(tx.Value, gasCost)

	// Check balance
	if senderAccount.Balance.Cmp(totalCost) < 0 {
		return &ExecutionResult{Status: 0, Error: ErrInsufficientBalance}, ErrInsufficientBalance
	}

	// Start execution
	gasUsed := uint64(21000) // Base gas cost
	logs := []*Log{}
	var contractAddress *crypto.Address

	// Deduct gas cost from sender
	senderAccount.Balance.Sub(senderAccount.Balance, gasCost)
	senderAccount.Nonce++

	if tx.IsContractCreation() {
		// Contract creation
		contractAddr := ee.generateContractAddress(tx.From, tx.Nonce-1)
		contractAddress = &contractAddr

		// Execute contract creation logic
		if len(tx.Data) > 0 {
			result, err := ee.executeContractCreation(tx, contractAddr, &gasUsed)
			if err != nil {
				// Refund remaining gas
				remainingGas := tx.GasLimit - gasUsed
				refund := new(big.Int).Mul(tx.GasPrice, big.NewInt(int64(remainingGas)))
				senderAccount.Balance.Add(senderAccount.Balance, refund)
				
				ee.stateDB.SetAccount(tx.From, senderAccount)
				return &ExecutionResult{
					GasUsed:         gasUsed,
					Status:          0,
					Logs:            logs,
					ContractAddress: contractAddress,
					Error:           err,
				}, nil
			}
			logs = append(logs, result.logs...)
		}

		// Create contract account
		contractAccount := &Account{
			Nonce:   1,
			Balance: new(big.Int).Set(tx.Value),
		}
		ee.stateDB.SetAccount(*contractAddress, contractAccount)
	} else {
		// Regular transfer or contract call
		if tx.To != nil {
			receiverAccount := ee.stateDB.GetAccount(*tx.To)
			if receiverAccount == nil {
				receiverAccount = &Account{
					Nonce:   0,
					Balance: big.NewInt(0),
				}
			}

			// Transfer value
			receiverAccount.Balance.Add(receiverAccount.Balance, tx.Value)
			ee.stateDB.SetAccount(*tx.To, receiverAccount)

			// Execute contract call if data is present
			if len(tx.Data) > 0 {
				result, err := ee.executeContractCall(tx, *tx.To, &gasUsed)
				if err != nil {
					// Contract call failed, but transaction succeeds
					// This is similar to Ethereum behavior
				}
				if result != nil {
					logs = append(logs, result.logs...)
				}
			}
		}
	}

	// Deduct value from sender
	senderAccount.Balance.Sub(senderAccount.Balance, tx.Value)

	// Refund remaining gas
	remainingGas := tx.GasLimit - gasUsed
	if remainingGas > 0 {
		refund := new(big.Int).Mul(tx.GasPrice, big.NewInt(int64(remainingGas)))
		senderAccount.Balance.Add(senderAccount.Balance, refund)
	}

	// Update sender account
	ee.stateDB.SetAccount(tx.From, senderAccount)

	return &ExecutionResult{
		GasUsed:         gasUsed,
		Status:          1,
		Logs:            logs,
		ContractAddress: contractAddress,
		Error:           nil,
	}, nil
}

// validateSignature validates the transaction signature
func (ee *ExecutionEngine) validateSignature(tx *Transaction) error {
	// Recreate transaction hash for signature verification
	hash := tx.CalculateHash()
	
	// Combine V, R, S into signature
	signature := make([]byte, 65)
	copy(signature[:32], tx.R.Bytes())
	copy(signature[32:64], tx.S.Bytes())
	signature[64] = byte(tx.V.Uint64())

	// Recover address from signature
	recoveredAddr, err := crypto.RecoverAddressFunc(hash, signature)
	if err != nil {
		return ErrInvalidSignature
	}

	// Check if recovered address matches the from address
	if !recoveredAddr.Equal(tx.From) {
		return ErrInvalidSignature
	}

	return nil
}

// generateContractAddress generates a contract address from sender and nonce
func (ee *ExecutionEngine) generateContractAddress(sender crypto.Address, nonce uint64) crypto.Address {
	// Simple implementation: hash(sender + nonce)
	data := append(sender.Bytes(), big.NewInt(int64(nonce)).Bytes()...)
	hash := crypto.BytesToHash(crypto.Keccak256(data))
	var addr crypto.Address
	copy(addr[:], hash[12:])
	return addr
}

// contractExecutionResult represents the result of contract execution
type contractExecutionResult struct {
	logs []*Log
}

// executeContractCreation executes contract creation logic
func (ee *ExecutionEngine) executeContractCreation(tx *Transaction, contractAddr crypto.Address, gasUsed *uint64) (*contractExecutionResult, error) {
	// Custom contract creation logic
	// For this implementation, we'll use a simple interpretation of the data field
	
	// Add base contract creation gas cost
	*gasUsed += 32000

	if *gasUsed > tx.GasLimit {
		return nil, ErrGasLimitExceeded
	}

	// Create a log for contract creation
	log := &Log{
		Address: contractAddr,
		Topics:  []crypto.Hash{crypto.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")}, // Contract created
		Data:    tx.Data,
	}

	return &contractExecutionResult{
		logs: []*Log{log},
	}, nil
}

// executeContractCall executes contract call logic
func (ee *ExecutionEngine) executeContractCall(tx *Transaction, contractAddr crypto.Address, gasUsed *uint64) (*contractExecutionResult, error) {
	// Custom contract call logic
	// For this implementation, we'll use a simple interpretation system
	
	// Add base contract call gas cost
	*gasUsed += 700

	if *gasUsed > tx.GasLimit {
		return nil, ErrGasLimitExceeded
	}

	// Simple contract execution based on data field
	if len(tx.Data) >= 4 {
		// Extract function selector (first 4 bytes)
		selector := tx.Data[:4]
		
		// Simple function implementations
		switch crypto.BytesToHash(selector).Hex() {
		case "0xa9059cbb": // transfer(address,uint256)
			return ee.executeTransfer(tx, contractAddr, gasUsed)
		case "0x70a08231": // balanceOf(address)
			return ee.executeBalanceOf(tx, contractAddr, gasUsed)
		default:
			// Unknown function, just create a generic log
			log := &Log{
				Address: contractAddr,
				Topics:  []crypto.Hash{crypto.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002")}, // Function called
				Data:    tx.Data,
			}
			return &contractExecutionResult{logs: []*Log{log}}, nil
		}
	}

	return &contractExecutionResult{logs: []*Log{}}, nil
}

// executeTransfer executes a token transfer function
func (ee *ExecutionEngine) executeTransfer(tx *Transaction, contractAddr crypto.Address, gasUsed *uint64) (*contractExecutionResult, error) {
	*gasUsed += 5000

	if *gasUsed > tx.GasLimit {
		return nil, ErrGasLimitExceeded
	}

	// Create transfer event log
	log := &Log{
		Address: contractAddr,
		Topics: []crypto.Hash{
			crypto.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"), // Transfer event
			crypto.BytesToHash(tx.From.Bytes()),                                                    // from
		},
		Data: tx.Data[4:], // Parameters
	}

	return &contractExecutionResult{logs: []*Log{log}}, nil
}

// executeBalanceOf executes a balance query function
func (ee *ExecutionEngine) executeBalanceOf(tx *Transaction, contractAddr crypto.Address, gasUsed *uint64) (*contractExecutionResult, error) {
	*gasUsed += 400

	if *gasUsed > tx.GasLimit {
		return nil, ErrGasLimitExceeded
	}

	// BalanceOf doesn't generate logs, just consumes gas
	return &contractExecutionResult{logs: []*Log{}}, nil
}

// EstimateGas estimates gas for a transaction
func (ee *ExecutionEngine) EstimateGas(tx *Transaction, header *BlockHeader) (uint64, error) {
	// Create a copy of the state for simulation
	stateDBCopy := ee.stateDB.Copy()
	engineCopy := &ExecutionEngine{
		stateDB: stateDBCopy,
		config:  ee.config,
	}

	// Simulate execution
	result, err := engineCopy.ExecuteTransaction(tx, header)
	if err != nil {
		return 0, err
	}

	// Add 10% buffer to the gas used
	estimatedGas := result.GasUsed * 11 / 10
	return estimatedGas, nil
}

// Call simulates a transaction call without state changes
func (ee *ExecutionEngine) Call(tx *Transaction, header *BlockHeader) ([]byte, error) {
	// Create a copy of the state for simulation
	stateDBCopy := ee.stateDB.Copy()
	engineCopy := &ExecutionEngine{
		stateDB: stateDBCopy,
		config:  ee.config,
	}

	// Simulate execution
	_, err := engineCopy.ExecuteTransaction(tx, header)
	if err != nil {
		return nil, err
	}

	// For this simple implementation, return empty data
	// In a real implementation, this would return the contract's return data
	return []byte{}, nil
}

// GetGasPrice returns the minimum gas price
func (ee *ExecutionEngine) GetGasPrice() *big.Int {
	return new(big.Int).Set(ee.config.MinGasPrice)
}

// GetChainID returns the chain ID
func (ee *ExecutionEngine) GetChainID() *big.Int {
	return new(big.Int).Set(ee.config.ChainID)
}
