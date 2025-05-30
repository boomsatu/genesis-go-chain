
package evm

import (
	"math/big"

	"blockchain-node/core"

	"github.com/ethereum/go-ethereum/common"
)

// StateDBAdapter adapts our StateDB to EVM interface requirements
type StateDBAdapter struct {
	stateDB *core.StateDB
}

// NewStateDBAdapter creates a new StateDB adapter
func NewStateDBAdapter(stateDB *core.StateDB) *StateDBAdapter {
	return &StateDBAdapter{
		stateDB: stateDB,
	}
}

// CreateAccount creates a new account
func (s *StateDBAdapter) CreateAccount(addr common.Address) {
	account := &core.Account{
		Nonce:   0,
		Balance: big.NewInt(0),
	}
	s.stateDB.SetAccount(addr, account)
}

// SubBalance subtracts amount from the account balance
func (s *StateDBAdapter) SubBalance(addr common.Address, amount *big.Int) {
	balance := s.stateDB.GetBalance(addr)
	newBalance := new(big.Int).Sub(balance, amount)
	s.stateDB.SetBalance(addr, newBalance)
}

// AddBalance adds amount to the account balance
func (s *StateDBAdapter) AddBalance(addr common.Address, amount *big.Int) {
	balance := s.stateDB.GetBalance(addr)
	newBalance := new(big.Int).Add(balance, amount)
	s.stateDB.SetBalance(addr, newBalance)
}

// GetBalance returns the account balance
func (s *StateDBAdapter) GetBalance(addr common.Address) *big.Int {
	return s.stateDB.GetBalance(addr)
}

// GetNonce returns the account nonce
func (s *StateDBAdapter) GetNonce(addr common.Address) uint64 {
	return s.stateDB.GetNonce(addr)
}

// SetNonce sets the account nonce
func (s *StateDBAdapter) SetNonce(addr common.Address, nonce uint64) {
	s.stateDB.SetNonce(addr, nonce)
}

// GetCodeHash returns the code hash of an account
func (s *StateDBAdapter) GetCodeHash(addr common.Address) common.Hash {
	account := s.stateDB.GetAccount(addr)
	if account == nil {
		return common.Hash{}
	}
	return account.CodeHash
}

// GetCode returns the code of an account
func (s *StateDBAdapter) GetCode(addr common.Address) []byte {
	return s.stateDB.GetCode(addr)
}

// SetCode sets the code for an account
func (s *StateDBAdapter) SetCode(addr common.Address, code []byte) {
	s.stateDB.SetCode(addr, code)
}

// GetCodeSize returns the size of the code
func (s *StateDBAdapter) GetCodeSize(addr common.Address) int {
	return len(s.stateDB.GetCode(addr))
}

// AddRefund adds to the refund counter
func (s *StateDBAdapter) AddRefund(gas uint64) {
	// Not implemented in our simple version
}

// SubRefund subtracts from the refund counter
func (s *StateDBAdapter) SubRefund(gas uint64) {
	// Not implemented in our simple version
}

// GetRefund returns the current refund counter
func (s *StateDBAdapter) GetRefund() uint64 {
	// Not implemented in our simple version
	return 0
}

// GetCommittedState returns the committed state value
func (s *StateDBAdapter) GetCommittedState(addr common.Address, key common.Hash) common.Hash {
	return s.stateDB.GetStorage(addr, key)
}

// GetState returns the current state value
func (s *StateDBAdapter) GetState(addr common.Address, key common.Hash) common.Hash {
	return s.stateDB.GetStorage(addr, key)
}

// SetState sets a state value
func (s *StateDBAdapter) SetState(addr common.Address, key, value common.Hash) {
	s.stateDB.SetStorage(addr, key, value)
}

// Suicide marks an account for deletion
func (s *StateDBAdapter) Suicide(addr common.Address) bool {
	// Simple implementation: just zero the balance
	s.stateDB.SetBalance(addr, big.NewInt(0))
	return true
}

// HasSuicided returns whether an account has been marked for deletion
func (s *StateDBAdapter) HasSuicided(addr common.Address) bool {
	// Simple implementation: check if balance is zero
	return s.stateDB.GetBalance(addr).Sign() == 0
}

// Exist checks if an account exists
func (s *StateDBAdapter) Exist(addr common.Address) bool {
	return s.stateDB.Exist(addr)
}

// Empty checks if an account is empty
func (s *StateDBAdapter) Empty(addr common.Address) bool {
	return s.stateDB.Empty(addr)
}

// RevertToSnapshot reverts state to a snapshot
func (s *StateDBAdapter) RevertToSnapshot(id int) {
	// Not implemented in our simple version
}

// Snapshot creates a state snapshot
func (s *StateDBAdapter) Snapshot() int {
	// Not implemented in our simple version
	return 0
}

// AddLog adds a log entry
func (s *StateDBAdapter) AddLog(log *core.Log) {
	s.stateDB.AddLog(log)
}

// AddPreimage adds a preimage to the cache
func (s *StateDBAdapter) AddPreimage(hash common.Hash, preimage []byte) {
	// Not implemented in our simple version
}

// ForEachStorage iterates over storage entries
func (s *StateDBAdapter) ForEachStorage(addr common.Address, cb func(key, value common.Hash) bool) error {
	// Not implemented in our simple version
	return nil
}
