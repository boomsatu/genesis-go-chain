
package evm

import (
	"math/big"

	"blockchain-node/core"
	"blockchain-node/crypto"
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
func (s *StateDBAdapter) CreateAccount(addr crypto.Address) {
	account := &core.Account{
		Nonce:   0,
		Balance: big.NewInt(0),
	}
	s.stateDB.SetAccount(addr, account)
}

// SubBalance subtracts amount from the account balance
func (s *StateDBAdapter) SubBalance(addr crypto.Address, amount *big.Int) {
	balance := s.stateDB.GetBalance(addr)
	newBalance := new(big.Int).Sub(balance, amount)
	s.stateDB.SetBalance(addr, newBalance)
}

// AddBalance adds amount to the account balance
func (s *StateDBAdapter) AddBalance(addr crypto.Address, amount *big.Int) {
	balance := s.stateDB.GetBalance(addr)
	newBalance := new(big.Int).Add(balance, amount)
	s.stateDB.SetBalance(addr, newBalance)
}

// GetBalance returns the account balance
func (s *StateDBAdapter) GetBalance(addr crypto.Address) *big.Int {
	return s.stateDB.GetBalance(addr)
}

// GetNonce returns the account nonce
func (s *StateDBAdapter) GetNonce(addr crypto.Address) uint64 {
	return s.stateDB.GetNonce(addr)
}

// SetNonce sets the account nonce
func (s *StateDBAdapter) SetNonce(addr crypto.Address, nonce uint64) {
	s.stateDB.SetNonce(addr, nonce)
}

// GetCodeHash returns the code hash of an account
func (s *StateDBAdapter) GetCodeHash(addr crypto.Address) crypto.Hash {
	account := s.stateDB.GetAccount(addr)
	if account == nil {
		return crypto.Hash{}
	}
	return account.CodeHash
}

// GetCode returns the code of an account
func (s *StateDBAdapter) GetCode(addr crypto.Address) []byte {
	return s.stateDB.GetCode(addr)
}

// SetCode sets the code for an account
func (s *StateDBAdapter) SetCode(addr crypto.Address, code []byte) {
	s.stateDB.SetCode(addr, code)
}

// GetCodeSize returns the size of the code
func (s *StateDBAdapter) GetCodeSize(addr crypto.Address) int {
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
func (s *StateDBAdapter) GetCommittedState(addr crypto.Address, key crypto.Hash) crypto.Hash {
	return s.stateDB.GetStorage(addr, key)
}

// GetState returns the current state value
func (s *StateDBAdapter) GetState(addr crypto.Address, key crypto.Hash) crypto.Hash {
	return s.stateDB.GetStorage(addr, key)
}

// SetState sets a state value
func (s *StateDBAdapter) SetState(addr crypto.Address, key, value crypto.Hash) {
	s.stateDB.SetStorage(addr, key, value)
}

// Suicide marks an account for deletion
func (s *StateDBAdapter) Suicide(addr crypto.Address) bool {
	// Simple implementation: just zero the balance
	s.stateDB.SetBalance(addr, big.NewInt(0))
	return true
}

// HasSuicided returns whether an account has been marked for deletion
func (s *StateDBAdapter) HasSuicided(addr crypto.Address) bool {
	// Simple implementation: check if balance is zero
	return s.stateDB.GetBalance(addr).Sign() == 0
}

// Exist checks if an account exists
func (s *StateDBAdapter) Exist(addr crypto.Address) bool {
	return s.stateDB.Exist(addr)
}

// Empty checks if an account is empty
func (s *StateDBAdapter) Empty(addr crypto.Address) bool {
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
func (s *StateDBAdapter) AddPreimage(hash crypto.Hash, preimage []byte) {
	// Not implemented in our simple version
}

// ForEachStorage iterates over storage entries
func (s *StateDBAdapter) ForEachStorage(addr crypto.Address, cb func(key, value crypto.Hash) bool) error {
	// Not implemented in our simple version
	return nil
}
