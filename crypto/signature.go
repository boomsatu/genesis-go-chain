
package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
)

// Sign calculates an ECDSA signature
func Sign(hash []byte, prv *ecdsa.PrivateKey) ([]byte, error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	
	r, s, err := ecdsa.Sign(rand.Reader, prv, hash)
	if err != nil {
		return nil, err
	}
	
	// Recovery ID calculation for Ethereum-style signatures
	recoveryId := 0
	for i := 0; i < 4; i++ {
		recoveredPub, err := recoverPublicKey(hash, r, s, i)
		if err != nil {
			continue
		}
		if recoveredPub.Equal(&prv.PublicKey) {
			recoveryId = i
			break
		}
	}
	
	// Encode signature: 32 bytes R + 32 bytes S + 1 byte recovery ID
	signature := make([]byte, 65)
	copy(signature[0:32], padBytes(r.Bytes(), 32))
	copy(signature[32:64], padBytes(s.Bytes(), 32))
	signature[64] = byte(recoveryId)
	
	return signature, nil
}

// VerifySignature checks that the given public key created signature over hash
func VerifySignature(pubkey, hash, signature []byte) bool {
	if len(signature) != 64 {
		return false
	}
	if len(hash) != 32 {
		return false
	}
	if len(pubkey) != 64 {
		return false
	}
	
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])
	
	x := new(big.Int).SetBytes(pubkey[:32])
	y := new(big.Int).SetBytes(pubkey[32:64])
	
	pub := &ecdsa.PublicKey{
		Curve: btcec.S256(),
		X:     x,
		Y:     y,
	}
	
	return ecdsa.Verify(pub, hash, r, s)
}

// SigToPub recovers the public key from a signature
func SigToPub(hash, sig []byte) (*ecdsa.PublicKey, error) {
	if len(sig) != 65 {
		return nil, fmt.Errorf("signature must be 65 bytes long")
	}
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash must be 32 bytes long")
	}
	
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])
	recoveryId := int(sig[64])
	
	return recoverPublicKey(hash, r, s, recoveryId)
}

// FromECDSA exports a private key into a binary dump
func FromECDSA(prv *ecdsa.PrivateKey) []byte {
	if prv == nil {
		return nil
	}
	return padBytes(prv.D.Bytes(), 32)
}

// ToECDSA creates a private key with the given D value
func ToECDSA(d []byte) (*ecdsa.PrivateKey, error) {
	return toECDSA(d, true)
}

// ToECDSAUnsafe blindly converts a binary blob to a private key
func ToECDSAUnsafe(d []byte) *ecdsa.PrivateKey {
	priv, _ := toECDSA(d, false)
	return priv
}

// FromECDSAPub exports a public key into a binary dump
func FromECDSAPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(btcec.S256(), pub.X, pub.Y)
}

// UnmarshalPubkey converts bytes to a secp256k1 public key
func UnmarshalPubkey(pub []byte) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(btcec.S256(), pub)
	if x == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return &ecdsa.PublicKey{Curve: btcec.S256(), X: x, Y: y}, nil
}

// HexToECDSA parses a secp256k1 private key
func HexToECDSA(hexkey string) (*ecdsa.PrivateKey, error) {
	b, err := FromHex(hexkey)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %v", err)
	}
	return ToECDSA(b)
}

// GenerateKey generates a new private key
func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(btcec.S256(), rand.Reader)
}

// ValidateSignatureValues verifies whether the signature values are valid
func ValidateSignatureValues(v byte, r, s *big.Int, homestead bool) bool {
	if r.Cmp(big.NewInt(1)) < 0 || s.Cmp(big.NewInt(1)) < 0 {
		return false
	}
	
	// Reject upper range of s values (ECDSA malleability)
	// See discussion in secp256k1/libsecp256k1/include/secp256k1.h
	if homestead && s.Cmp(secp256k1halfN) > 0 {
		return false
	}
	
	// Frontier: allow s to be in full N range
	return r.Cmp(secp256k1N) < 0 && s.Cmp(secp256k1N) < 0 && (v == 0 || v == 1)
}

// Helper functions

var (
	secp256k1N, _     = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN   = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

func toECDSA(d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = btcec.S256()
	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)
	
	// The priv.D must < N
	if priv.D.Cmp(secp256k1N) >= 0 {
		return nil, fmt.Errorf("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, fmt.Errorf("invalid private key, zero or negative")
	}
	
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, fmt.Errorf("invalid private key")
	}
	return priv, nil
}

func recoverPublicKey(hash []byte, r, s *big.Int, recoveryId int) (*ecdsa.PublicKey, error) {
	curve := btcec.S256()
	
	// Calculate point R = (r, y)
	x := r
	if recoveryId >= 2 {
		x = new(big.Int).Add(r, curve.Params().N)
	}
	
	// Calculate y coordinate
	y2 := new(big.Int).Exp(x, big.NewInt(3), curve.Params().P)
	y2.Add(y2, curve.Params().B)
	y := new(big.Int).ModSqrt(y2, curve.Params().P)
	
	if y == nil {
		return nil, fmt.Errorf("invalid recovery id")
	}
	
	if y.Bit(0) != uint(recoveryId&1) {
		y.Sub(curve.Params().P, y)
	}
	
	// Calculate e = -hash mod N
	e := new(big.Int).SetBytes(hash)
	e.Neg(e)
	e.Mod(e, curve.Params().N)
	
	// Calculate r^-1 mod N
	rInv := new(big.Int).ModInverse(r, curve.Params().N)
	
	// Calculate point Q = r^-1 * (s*R - e*G)
	sR_x, sR_y := curve.ScalarMult(x, y, s.Bytes())
	eG_x, eG_y := curve.ScalarBaseMult(e.Bytes())
	
	// Subtract eG from sR
	eG_y.Neg(eG_y)
	eG_y.Mod(eG_y, curve.Params().P)
	
	Q_x, Q_y := curve.Add(sR_x, sR_y, eG_x, eG_y)
	
	// Multiply by r^-1
	Q_x, Q_y = curve.ScalarMult(Q_x, Q_y, rInv.Bytes())
	
	return &ecdsa.PublicKey{
		Curve: curve,
		X:     Q_x,
		Y:     Q_y,
	}, nil
}

func padBytes(b []byte, size int) []byte {
	if len(b) >= size {
		return b
	}
	padded := make([]byte, size)
	copy(padded[size-len(b):], b)
	return padded
}
