
package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
)

// Wallet represents a cryptocurrency wallet
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
	Address    Address
}

// NewWallet creates a new wallet with a random private key
func NewWallet() (*Wallet, error) {
	privateKey, err := ecdsa.GenerateKey(btcec.S256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	publicKey := FromECDSAPub(&privateKey.PublicKey)
	address := PubkeyToAddress(publicKey)

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}, nil
}

// WalletFromPrivateKey creates a wallet from an existing private key
func WalletFromPrivateKey(privateKey *ecdsa.PrivateKey) *Wallet {
	publicKey := FromECDSAPub(&privateKey.PublicKey)
	address := PubkeyToAddress(publicKey)

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}
}

// GetPrivateKeyHex returns the private key as a hex string
func (w *Wallet) GetPrivateKeyHex() string {
	return fmt.Sprintf("%x", FromECDSA(w.PrivateKey))
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
func (w *Wallet) SignHash(hash Hash) ([]byte, error) {
	signature, err := Sign(hash.Bytes(), w.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign hash: %v", err)
	}
	return signature, nil
}

// VerifySignatureFunc verifies a signature against a hash and public key
func VerifySignatureFunc(hash Hash, signature []byte, publicKey []byte) bool {
	return VerifySignature(publicKey, hash.Bytes(), signature[:64])
}

// RecoverPublicKeyFunc recovers the public key from a signature and hash
func RecoverPublicKeyFunc(hash Hash, signature []byte) ([]byte, error) {
	publicKey, err := SigToPub(hash.Bytes(), signature)
	if err != nil {
		return nil, fmt.Errorf("failed to recover public key: %v", err)
	}
	return FromECDSAPub(publicKey), nil
}

// RecoverAddressFunc recovers the address from a signature and hash
func RecoverAddressFunc(hash Hash, signature []byte) (Address, error) {
	publicKey, err := SigToPub(hash.Bytes(), signature)
	if err != nil {
		return Address{}, fmt.Errorf("failed to recover address: %v", err)
	}
	return PubkeyToAddress(FromECDSAPub(publicKey)), nil
}
