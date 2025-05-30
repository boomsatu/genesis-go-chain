
package crypto

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	// AddressLength is the expected length of the address
	AddressLength = 20
	// HashLength is the expected length of the hash
	HashLength = 32
)

// Address represents a 20 byte address of an Ethereum account
type Address [AddressLength]byte

// Hash represents a 32 byte Keccak256 hash
type Hash [HashLength]byte

// BytesToAddress returns Address with value b
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

// HexToAddress returns Address with byte values of s
func HexToAddress(s string) Address {
	return BytesToAddress(FromHex(s))
}

// IsHexAddress verifies whether a string can represent a valid hex-encoded address or not
func IsHexAddress(s string) bool {
	if has0xPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && isHex(s)
}

// SetBytes sets the address to the value of b
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// Bytes gets the string representation of the underlying address
func (a Address) Bytes() []byte {
	return a[:]
}

// Hex returns a hex string representation of the address
func (a Address) Hex() string {
	return string(a.checksumHex())
}

// String implements fmt.Stringer
func (a Address) String() string {
	return a.Hex()
}

// checksumHex returns the hex representation of the address with EIP-55 checksum
func (a Address) checksumHex() []byte {
	buf := a.hex()
	hash := Keccak256(buf[2:])
	for i := 2; i < len(buf); i++ {
		hashByte := hash[(i-2)/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if buf[i] > '9' && hashByte > 7 {
			buf[i] -= 32
		}
	}
	return buf[:]
}

// hex returns the hex representation of the address
func (a Address) hex() []byte {
	var buf [len(a)*2 + 2]byte
	copy(buf[:2], "0x")
	hex.Encode(buf[2:], a[:])
	return buf[:]
}

// BytesToHash sets b to hash
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// HexToHash sets byte representation of s to hash
func HexToHash(s string) Hash {
	return BytesToHash(FromHex(s))
}

// SetBytes sets the hash to the value of b
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}
	copy(h[HashLength-len(b):], b)
}

// Bytes gets the byte representation of the underlying hash
func (h Hash) Bytes() []byte {
	return h[:]
}

// Hex converts a hash to a hex string
func (h Hash) Hex() string {
	return hexEncodeToString(h[:])
}

// String implements the fmt.Stringer interface
func (h Hash) String() string {
	return h.Hex()
}

// FromHex returns the bytes represented by the hexadecimal string s
func FromHex(s string) []byte {
	if has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// has0xPrefix validates str begins with '0x' or '0X'
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// isHex validates whether each byte is valid hexadecimal string
func isHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

// isHexCharacter returns bool of c being a valid hexadecimal
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

// hexEncodeToString encodes b as a hex string with 0x prefix
func hexEncodeToString(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
}

// Equal returns true if both addresses are equal
func (a Address) Equal(other Address) bool {
	return bytes.Equal(a[:], other[:])
}

// Equal returns true if both hashes are equal
func (h Hash) Equal(other Hash) bool {
	return bytes.Equal(h[:], other[:])
}

// EmptyAddress returns an empty address
func EmptyAddress() Address {
	return Address{}
}

// EmptyHash returns an empty hash
func EmptyHash() Hash {
	return Hash{}
}

// PubkeyToAddress creates an address from a public key
func PubkeyToAddress(pubkey []byte) Address {
	if len(pubkey) == 65 {
		pubkey = pubkey[1:] // Remove the 0x04 prefix for uncompressed key
	}
	return BytesToAddress(Keccak256(pubkey)[12:])
}

// AddressFromString creates an address from a hex string
func AddressFromString(s string) (Address, error) {
	if !IsHexAddress(s) {
		return Address{}, fmt.Errorf("invalid address format: %s", s)
	}
	return HexToAddress(s), nil
}

// HashFromString creates a hash from a hex string
func HashFromString(s string) (Hash, error) {
	if !has0xPrefix(s) {
		s = "0x" + s
	}
	if len(s) != 2+2*HashLength {
		return Hash{}, fmt.Errorf("invalid hash length: expected %d, got %d", 2+2*HashLength, len(s))
	}
	if !isHex(s[2:]) {
		return Hash{}, fmt.Errorf("invalid hex string: %s", s)
	}
	return HexToHash(s), nil
}

// MustAddressFromString creates an address from a hex string, panics on error
func MustAddressFromString(s string) Address {
	addr, err := AddressFromString(s)
	if err != nil {
		panic(err)
	}
	return addr
}

// MustHashFromString creates a hash from a hex string, panics on error
func MustHashFromString(s string) Hash {
	hash, err := HashFromString(s)
	if err != nil {
		panic(err)
	}
	return hash
}

// ToLower returns the address in lowercase
func (a Address) ToLower() string {
	return strings.ToLower(a.Hex())
}

// IsZero returns true if the address is zero
func (a Address) IsZero() bool {
	return bytes.Equal(a[:], EmptyAddress()[:])
}

// IsZero returns true if the hash is zero
func (h Hash) IsZero() bool {
	return bytes.Equal(h[:], EmptyHash()[:])
}
