
package evm

import (
	"math/big"

	"blockchain-node/core"
	"blockchain-node/storage"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
)

// StateDB implements vm.StateDB interface for EVM integration
type StateDB struct {
	db       storage.Database
	accounts map[common.Address]*core.Account
	storage  map[common.Address]map[common.Hash]common.Hash
	logs     []*core.Log
	
	// Transaction context
	txHash      common.Hash
	txIndex     int
	blockHash   common.Hash
	blockNumber *big.Int
	
	// Snapshots for rollback
	snapshots   []snapshot
	snapshotId  int
}

type snapshot struct {
	id       int
	accounts map[common.Address]*core.Account
	storage  map[common.Address]map[common.Hash]common.Hash
}

// NewStateDB creates a new StateDB instance
func NewStateDB(db storage.Database) *StateDB {
	return &StateDB{
		db:       db,
		accounts: make(map[common.Address]*core.Account),
		storage:  make(map[common.Address]map[common.Hash]common.Hash),
		logs:     make([]*core.Log, 0),
	}
}

// CreateAccount creates a new account
func (s *StateDB) CreateAccount(addr common.Address) {
	s.accounts[addr] = &core.Account{
		Nonce:       0,
		Balance:     big.NewInt(0),
		CodeHash:    crypto.Keccak256Hash(nil),
		StorageRoot: common.Hash{},
	}
}

// SubBalance subtracts amount from the account associated with addr
func (s *StateDB) SubBalance(addr common.Address, amount *big.Int) {
	account := s.getAccount(addr)
	account.Balance = new(big.Int).Sub(account.Balance, amount)
}

// AddBalance adds amount to the account associated with addr
func (s *StateDB) AddBalance(addr common.Address, amount *big.Int) {
	account := s.getAccount(addr)
	account.Balance = new(big.Int).Add(account.Balance, amount)
}

// GetBalance retrieves the balance from the given address
func (s *StateDB) GetBalance(addr common.Address) *big.Int {
	account := s.getAccount(addr)
	return new(big.Int).Set(account.Balance)
}

// GetNonce returns the nonce of the given address
func (s *StateDB) GetNonce(addr common.Address) uint64 {
	account := s.getAccount(addr)
	return account.Nonce
}

// SetNonce sets the nonce of the given address
func (s *StateDB) SetNonce(addr common.Address, nonce uint64) {
	account := s.getAccount(addr)
	account.Nonce = nonce
}

// GetCodeHash returns the code hash of the given address
func (s *StateDB) GetCodeHash(addr common.Address) common.Hash {
	account := s.getAccount(addr)
	return account.CodeHash
}

// GetCode returns the code associated with the given address
func (s *StateDB) GetCode(addr common.Address) []byte {
	codeHash := s.GetCodeHash(addr)
	if codeHash == crypto.Keccak256Hash(nil) {
		return nil
	}
	
	// Retrieve code from database
	code, err := s.db.Get(append([]byte("code-"), codeHash.Bytes()...))
	if err != nil {
		return nil
	}
	return code
}

// SetCode sets the code associated with the given address
func (s *StateDB) SetCode(addr common.Address, code []byte) {
	codeHash := crypto.Keccak256Hash(code)
	account := s.getAccount(addr)
	account.CodeHash = codeHash
	
	// Store code in database
	s.db.Put(append([]byte("code-"), codeHash.Bytes()...), code)
}

// GetCodeSize returns the size of the code associated with the given address
func (s *StateDB) GetCodeSize(addr common.Address) int {
	code := s.GetCode(addr)
	return len(code)
}

// AddRefund adds gas to the refund counter
func (s *StateDB) AddRefund(gas uint64) {
	// TODO: Implement refund mechanism
}

// SubRefund subtracts gas from the refund counter
func (s *StateDB) SubRefund(gas uint64) {
	// TODO: Implement refund mechanism
}

// GetRefund returns the current value of the refund counter
func (s *StateDB) GetRefund() uint64 {
	// TODO: Implement refund mechanism
	return 0
}

// GetCommittedState retrieves a value from the given account's committed storage trie
func (s *StateDB) GetCommittedState(addr common.Address, key common.Hash) common.Hash {
	return s.GetState(addr, key)
}

// GetState retrieves a value from the given account's storage trie
func (s *StateDB) GetState(addr common.Address, key common.Hash) common.Hash {
	if s.storage[addr] == nil {
		s.storage[addr] = make(map[common.Hash]common.Hash)
		// TODO: Load from database/trie
	}
	return s.storage[addr][key]
}

// SetState sets a value in the given account's storage trie
func (s *StateDB) SetState(addr common.Address, key, value common.Hash) {
	if s.storage[addr] == nil {
		s.storage[addr] = make(map[common.Hash]common.Hash)
	}
	s.storage[addr][key] = value
}

// GetTransientState gets transient storage for a given account
func (s *StateDB) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	// TODO: Implement transient storage
	return common.Hash{}
}

// SetTransientState sets transient storage for a given account
func (s *StateDB) SetTransientState(addr common.Address, key, value common.Hash) {
	// TODO: Implement transient storage
}

// Suicide marks the given account as suicided
func (s *StateDB) Suicide(addr common.Address) bool {
	// TODO: Implement suicide mechanism
	return false
}

// HasSuicided returns if the account is suicided
func (s *StateDB) HasSuicided(addr common.Address) bool {
	// TODO: Implement suicide check
	return false
}

// Exist reports whether the given account exists in state
func (s *StateDB) Exist(addr common.Address) bool {
	_, exists := s.accounts[addr]
	return exists
}

// Empty returns whether the given account is empty
func (s *StateDB) Empty(addr common.Address) bool {
	account := s.getAccount(addr)
	return account.Nonce == 0 && account.Balance.Sign() == 0 && account.CodeHash == crypto.Keccak256Hash(nil)
}

// AddressInAccessList returns true if the address is in the access list
func (s *StateDB) AddressInAccessList(addr common.Address) bool {
	// TODO: Implement access list
	return false
}

// SlotInAccessList returns true if the (address, slot)-tuple is in the access list
func (s *StateDB) SlotInAccessList(addr common.Address, slot common.Hash) (addressPresent bool, slotPresent bool) {
	// TODO: Implement access list
	return false, false
}

// AddAddressToAccessList adds the given address to the access list
func (s *StateDB) AddAddressToAccessList(addr common.Address) {
	// TODO: Implement access list
}

// AddSlotToAccessList adds the given (address, slot)-tuple to the access list
func (s *StateDB) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	// TODO: Implement access list
}

// RevertToSnapshot reverts all state changes made since the given revision
func (s *StateDB) RevertToSnapshot(revid int) {
	// Find the snapshot
	for i := len(s.snapshots) - 1; i >= 0; i-- {
		if s.snapshots[i].id == revid {
			// Restore state
			s.accounts = make(map[common.Address]*core.Account)
			for addr, acc := range s.snapshots[i].accounts {
				s.accounts[addr] = &core.Account{
					Nonce:       acc.Nonce,
					Balance:     new(big.Int).Set(acc.Balance),
					CodeHash:    acc.CodeHash,
					StorageRoot: acc.StorageRoot,
				}
			}
			s.storage = make(map[common.Address]map[common.Hash]common.Hash)
			for addr, storage := range s.snapshots[i].storage {
				s.storage[addr] = make(map[common.Hash]common.Hash)
				for key, value := range storage {
					s.storage[addr][key] = value
				}
			}
			// Remove snapshots after this one
			s.snapshots = s.snapshots[:i]
			break
		}
	}
}

// Snapshot returns an identifier for the current revision of the state
func (s *StateDB) Snapshot() int {
	id := s.snapshotId
	s.snapshotId++
	
	// Create deep copy of current state
	accountsCopy := make(map[common.Address]*core.Account)
	for addr, acc := range s.accounts {
		accountsCopy[addr] = &core.Account{
			Nonce:       acc.Nonce,
			Balance:     new(big.Int).Set(acc.Balance),
			CodeHash:    acc.CodeHash,
			StorageRoot: acc.StorageRoot,
		}
	}
	
	storageCopy := make(map[common.Address]map[common.Hash]common.Hash)
	for addr, storage := range s.storage {
		storageCopy[addr] = make(map[common.Hash]common.Hash)
		for key, value := range storage {
			storageCopy[addr][key] = value
		}
	}
	
	s.snapshots = append(s.snapshots, snapshot{
		id:       id,
		accounts: accountsCopy,
		storage:  storageCopy,
	})
	
	return id
}

// AddLog adds a log
func (s *StateDB) AddLog(log *types.Log) {
	s.logs = append(s.logs, &core.Log{
		Address:     log.Address,
		Topics:      log.Topics,
		Data:        log.Data,
		BlockNumber: s.blockNumber.Uint64(),
		TxHash:      s.txHash,
		TxIndex:     uint(s.txIndex),
		BlockHash:   s.blockHash,
		Index:       uint(len(s.logs)),
		Removed:     false,
	})
}

// GetLogs returns all logs
func (s *StateDB) GetLogs() []*core.Log {
	return s.logs
}

// Prepare sets the current transaction context
func (s *StateDB) Prepare(txHash common.Hash, txIndex int) {
	s.txHash = txHash
	s.txIndex = txIndex
}

// SetBlockContext sets the current block context
func (s *StateDB) SetBlockContext(blockHash common.Hash, blockNumber *big.Int) {
	s.blockHash = blockHash
	s.blockNumber = blockNumber
}

// getAccount gets or creates an account
func (s *StateDB) getAccount(addr common.Address) *core.Account {
	if account, exists := s.accounts[addr]; exists {
		return account
	}
	
	// Try to load from database
	// TODO: Implement loading from database
	
	// Create new account if not found
	s.CreateAccount(addr)
	return s.accounts[addr]
}
