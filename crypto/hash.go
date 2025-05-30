
package crypto

import (
	"crypto/sha256"
	"hash"
	
	"golang.org/x/crypto/sha3"
)

// Keccak256 calculates and returns the Keccak256 hash of the input data
func Keccak256(data ...[]byte) []byte {
	hasher := sha3.NewLegacyKeccak256()
	for _, b := range data {
		hasher.Write(b)
	}
	return hasher.Sum(nil)
}

// Keccak256Hash calculates and returns the Keccak256 hash as a Hash
func Keccak256Hash(data ...[]byte) Hash {
	return BytesToHash(Keccak256(data...))
}

// Sha256 calculates and returns the SHA256 hash of the input data
func Sha256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// NewKeccak256 creates a new Keccak256 hasher
func NewKeccak256() hash.Hash {
	return sha3.NewLegacyKeccak256()
}
