
package core

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"

	"blockchain-node/crypto"
	"blockchain-node/storage"
)

// StateDB manages the world state using Patricia Merkle Trie structure
type StateDB struct {
	db       storage.Database
	stateRoot crypto.Hash
	accounts  map[crypto.Address]*Account // In-memory cache
	storage   map[crypto.Address]map[crypto.Hash]crypto.Hash // Contract storage
	logs      []*Log
	mu        sync.RWMutex
}

// NewStateDB creates a new StateDB instance
func NewStateDB(db storage.Database, stateRoot crypto.Hash) *StateDB {
	return &StateDB{
		db:        db,
		stateRoot: stateRoot,
		accounts:  make(map[crypto.Address]*Account),
		storage:   make(map[crypto.Address]map[crypto.Hash]crypto.Hash),
		logs:      []*Log{},
	}
}

// GetAccount retrieves an account from the state
func (sdb *StateDB) GetAccount(addr crypto.Address) *Account {
	sdb.mu.RLock()
	defer sdb.mu.RUnlock()

	// Check cache first
	if account, exists := sdb.accounts[addr]; exists {
		return account
	}

	// Load from database
	key := append([]byte("account-"), addr.Bytes()...)
	data, err := sdb.db.Get(key)
	if err != nil {
		return nil
	}

	var account Account
	if err := json.Unmarshal(data, &account); err != nil {
		return nil
	}

	// Cache the account
	sdb.accounts[addr] = &account
	return &account
}

// SetAccount updates an account in the state
func (sdb *StateDB) SetAccount(addr crypto.Address, account *Account) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	// Update cache
	sdb.accounts[addr] = account
}

// GetBalance returns the balance of an account
func (sdb *StateDB) GetBalance(addr crypto.Address) *big.Int {
	account := sdb.GetAccount(addr)
	if account == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(account.Balance)
}

// SetBalance updates the balance of an account
func (sdb *StateDB) SetBalance(addr crypto.Address, balance *big.Int) {
	account := sdb.GetAccount(addr)
	if account == nil {
		account = &Account{
			Nonce:   0,
			Balance: big.NewInt(0),
		}
	}
	account.Balance = new(big.Int).Set(balance)
	sdb.SetAccount(addr, account)
}

// GetNonce returns the nonce of an account
func (sdb *StateDB) GetNonce(addr crypto.Address) uint64 {
	account := sdb.GetAccount(addr)
	if account == nil {
		return 0
	}
	return account.Nonce
}

// SetNonce updates the nonce of an account
func (sdb *StateDB) SetNonce(addr crypto.Address, nonce uint64) {
	account := sdb.GetAccount(addr)
	if account == nil {
		account = &Account{
			Nonce:   0,
			Balance: big.NewInt(0),
		}
	}
	account.Nonce = nonce
	sdb.SetAccount(addr, account)
}

// GetCode returns the code of a contract account
func (sdb *StateDB) GetCode(addr crypto.Address) []byte {
	account := sdb.GetAccount(addr)
	if account == nil {
		return nil
	}

	if account.CodeHash.IsZero() {
		return nil
	}

	// Load code from database
	key := append([]byte("code-"), account.CodeHash.Bytes()...)
	data, err := sdb.db.Get(key)
	if err != nil {
		return nil
	}

	return data
}

// SetCode updates the code of a contract account
func (sdb *StateDB) SetCode(addr crypto.Address, code []byte) {
	account := sdb.GetAccount(addr)
	if account == nil {
		account = &Account{
			Nonce:   0,
			Balance: big.NewInt(0),
		}
	}

	// Calculate code hash
	codeHash := crypto.Keccak256Hash(code)
	account.CodeHash = codeHash

	// Store code in database
	key := append([]byte("code-"), codeHash.Bytes()...)
	sdb.db.Put(key, code)

	sdb.SetAccount(addr, account)
}

// GetStorage returns a storage value for a contract
func (sdb *StateDB) GetStorage(addr crypto.Address, key crypto.Hash) crypto.Hash {
	sdb.mu.RLock()
	defer sdb.mu.RUnlock()

	// Check cache first
	if addrStorage, exists := sdb.storage[addr]; exists {
		if value, exists := addrStorage[key]; exists {
			return value
		}
	}

	// Load from database
	dbKey := append([]byte("storage-"), addr.Bytes()...)
	dbKey = append(dbKey, key.Bytes()...)
	
	data, err := sdb.db.Get(dbKey)
	if err != nil {
		return crypto.Hash{}
	}

	value := crypto.BytesToHash(data)

	// Cache the value
	if sdb.storage[addr] == nil {
		sdb.storage[addr] = make(map[crypto.Hash]crypto.Hash)
	}
	sdb.storage[addr][key] = value

	return value
}

// SetStorage updates a storage value for a contract
func (sdb *StateDB) SetStorage(addr crypto.Address, key crypto.Hash, value crypto.Hash) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	// Update cache
	if sdb.storage[addr] == nil {
		sdb.storage[addr] = make(map[crypto.Hash]crypto.Hash)
	}
	sdb.storage[addr][key] = value
}

// AddLog adds a log to the state
func (sdb *StateDB) AddLog(log *Log) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()
	
	sdb.logs = append(sdb.logs, log)
}

// GetLogs returns all logs in the current state
func (sdb *StateDB) GetLogs() []*Log {
	sdb.mu.RLock()
	defer sdb.mu.RUnlock()
	
	return append([]*Log{}, sdb.logs...)
}

// Commit commits all changes to the database and returns the new state root
func (sdb *StateDB) Commit() (crypto.Hash, error) {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	// Create a batch for atomic writes
	batch := sdb.db.NewBatch()

	// Commit all account changes
	for addr, account := range sdb.accounts {
		data, err := json.Marshal(account)
		if err != nil {
			return crypto.Hash{}, fmt.Errorf("failed to marshal account: %v", err)
		}

		key := append([]byte("account-"), addr.Bytes()...)
		if err := batch.Put(key, data); err != nil {
			return crypto.Hash{}, fmt.Errorf("failed to put account: %v", err)
		}
	}

	// Commit all storage changes
	for addr, addrStorage := range sdb.storage {
		for key, value := range addrStorage {
			dbKey := append([]byte("storage-"), addr.Bytes()...)
			dbKey = append(dbKey, key.Bytes()...)
			
			if err := batch.Put(dbKey, value.Bytes()); err != nil {
				return crypto.Hash{}, fmt.Errorf("failed to put storage: %v", err)
			}
		}
	}

	// Write the batch
	if err := batch.Write(); err != nil {
		return crypto.Hash{}, fmt.Errorf("failed to write batch: %v", err)
	}

	// Calculate new state root
	newStateRoot := sdb.calculateStateRoot()
	sdb.stateRoot = newStateRoot

	// Clear caches
	sdb.accounts = make(map[crypto.Address]*Account)
	sdb.storage = make(map[crypto.Address]map[crypto.Hash]crypto.Hash)
	sdb.logs = []*Log{}

	return newStateRoot, nil
}

// calculateStateRoot calculates the state root using a simple merkle tree
func (sdb *StateDB) calculateStateRoot() crypto.Hash {
	// Simple implementation: hash all account addresses and balances
	// In a real implementation, this would be a proper Patricia Merkle Trie
	
	var data []byte
	
	// Add accounts to hash calculation
	for addr, account := range sdb.accounts {
		data = append(data, addr.Bytes()...)
		data = append(data, account.Balance.Bytes()...)
		data = append(data, big.NewInt(int64(account.Nonce)).Bytes()...)
		data = append(data, account.CodeHash.Bytes()...)
	}

	// Add storage to hash calculation
	for addr, addrStorage := range sdb.storage {
		data = append(data, addr.Bytes()...)
		for key, value := range addrStorage {
			data = append(data, key.Bytes()...)
			data = append(data, value.Bytes()...)
		}
	}

	if len(data) == 0 {
		return crypto.Hash{}
	}

	return crypto.Keccak256Hash(data)
}

// Copy creates a deep copy of the StateDB
func (sdb *StateDB) Copy() *StateDB {
	sdb.mu.RLock()
	defer sdb.mu.RUnlock()

	copy := &StateDB{
		db:        sdb.db,
		stateRoot: sdb.stateRoot,
		accounts:  make(map[crypto.Address]*Account),
		storage:   make(map[crypto.Address]map[crypto.Hash]crypto.Hash),
		logs:      make([]*Log, len(sdb.logs)),
	}

	// Copy accounts
	for addr, account := range sdb.accounts {
		copy.accounts[addr] = &Account{
			Nonce:       account.Nonce,
			Balance:     new(big.Int).Set(account.Balance),
			CodeHash:    account.CodeHash,
			StorageRoot: account.StorageRoot,
		}
	}

	// Copy storage
	for addr, addrStorage := range sdb.storage {
		copy.storage[addr] = make(map[crypto.Hash]crypto.Hash)
		for key, value := range addrStorage {
			copy.storage[addr][key] = value
		}
	}

	// Copy logs
	for i, log := range sdb.logs {
		copy.logs[i] = &Log{
			Address:     log.Address,
			Topics:      append([]crypto.Hash{}, log.Topics...),
			Data:        append([]byte{}, log.Data...),
			BlockNumber: log.BlockNumber,
			TxHash:      log.TxHash,
			TxIndex:     log.TxIndex,
			BlockHash:   log.BlockHash,
			Index:       log.Index,
			Removed:     log.Removed,
		}
	}

	return copy
}

// GetStateRoot returns the current state root
func (sdb *StateDB) GetStateRoot() crypto.Hash {
	sdb.mu.RLock()
	defer sdb.mu.RUnlock()
	return sdb.stateRoot
}

// Empty checks if an account is empty (non-existent or with zero nonce, balance, and no code)
func (sdb *StateDB) Empty(addr crypto.Address) bool {
	account := sdb.GetAccount(addr)
	if account == nil {
		return true
	}

	return account.Nonce == 0 &&
		account.Balance.Sign() == 0 &&
		account.CodeHash.IsZero()
}

// Exist checks if an account exists in the state
func (sdb *StateDB) Exist(addr crypto.Address) bool {
	return sdb.GetAccount(addr) != nil
}

// GetAccountsCount returns the number of accounts in the cache
func (sdb *StateDB) GetAccountsCount() int {
	sdb.mu.RLock()
	defer sdb.mu.RUnlock()
	return len(sdb.accounts)
}
