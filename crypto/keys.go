
package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Wallet represents a cryptocurrency wallet
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
	Address    common.Address
}

// NewWallet creates a new wallet with a random private key
func NewWallet() (*Wallet, error) {
	privateKey, err := ecdsa.GenerateKey(btcec.S256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	publicKey := crypto.FromECDSAPub(&privateKey.PublicKey)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}, nil
}

// WalletFromPrivateKey creates a wallet from an existing private key
func WalletFromPrivateKey(privateKey *ecdsa.PrivateKey) *Wallet {
	publicKey := crypto.FromECDSAPub(&privateKey.PublicKey)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}
}

// GetPrivateKeyHex returns the private key as a hex string
func (w *Wallet) GetPrivateKeyHex() string {
	return fmt.Sprintf("%x", crypto.FromECDSA(w.PrivateKey))
}

// GetPublicKeyHex returns the public key as a hex string
func (w *Wallet) GetPublicKeyHex() string {
	return fmt.Sprintf("%x", w.PublicKey)
}

// GetAddressHex returns the address as a hex string
func (w *Wallet) GetAddressHex() string {
	return w.Address.Hex()
}

// SignHash signs a hash with the wallet's private key
func (w *Wallet) SignHash(hash common.Hash) ([]byte, error) {
	signature, err := crypto.Sign(hash.Bytes(), w.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign hash: %v", err)
	}
	return signature, nil
}

// VerifySignature verifies a signature against a hash and public key
func VerifySignature(hash common.Hash, signature []byte, publicKey []byte) bool {
	return crypto.VerifySignature(publicKey, hash.Bytes(), signature[:64])
}

// RecoverPublicKey recovers the public key from a signature and hash
func RecoverPublicKey(hash common.Hash, signature []byte) ([]byte, error) {
	publicKey, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		return nil, fmt.Errorf("failed to recover public key: %v", err)
	}
	return crypto.FromECDSAPub(publicKey), nil
}

// RecoverAddress recovers the address from a signature and hash
func RecoverAddress(hash common.Hash, signature []byte) (common.Address, error) {
	publicKey, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to recover address: %v", err)
	}
	return crypto.PubkeyToAddress(*publicKey), nil
}
