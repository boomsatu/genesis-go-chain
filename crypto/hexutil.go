
package crypto

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
)

// EncodeUint64 encodes i as a hex string with 0x prefix
func EncodeUint64(i uint64) string {
	enc := strconv.FormatUint(i, 16)
	if len(enc)%2 != 0 {
		enc = "0" + enc
	}
	return "0x" + enc
}

// DecodeUint64 decodes a hex string with 0x prefix as uint64
func DecodeUint64(input string) (uint64, error) {
	if !has0xPrefix(input) {
		return 0, fmt.Errorf("hex string without 0x prefix")
	}
	raw := input[2:]
	if len(raw) == 0 {
		return 0, fmt.Errorf("empty hex string")
	}
	return strconv.ParseUint(raw, 16, 64)
}

// EncodeBig encodes bigint as a hex string with 0x prefix
func EncodeBig(bigint *big.Int) string {
	if bigint == nil {
		return "0x0"
	}
	if bigint.Sign() == 0 {
		return "0x0"
	}
	return "0x" + bigint.Text(16)
}

// DecodeBig decodes a hex string with 0x prefix as *big.Int
func DecodeBig(input string) (*big.Int, error) {
	if !has0xPrefix(input) {
		return nil, fmt.Errorf("hex string without 0x prefix")
	}
	raw := input[2:]
	if len(raw) == 0 {
		return nil, fmt.Errorf("empty hex string")
	}
	bigint, ok := new(big.Int).SetString(raw, 16)
	if !ok {
		return nil, fmt.Errorf("invalid hex string")
	}
	return bigint, nil
}

// Encode encodes b as a hex string with 0x prefix
func Encode(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

// Decode decodes a hex string with 0x prefix
func Decode(input string) ([]byte, error) {
	if !has0xPrefix(input) {
		return nil, fmt.Errorf("hex string without 0x prefix")
	}
	raw := input[2:]
	if len(raw)%2 != 0 {
		raw = "0" + raw
	}
	return hex.DecodeString(raw)
}

// MustDecode decodes a hex string with 0x prefix, panics on error
func MustDecode(input string) []byte {
	dec, err := Decode(input)
	if err != nil {
		panic(err)
	}
	return dec
}

// EncodeToString encodes b as a hex string with 0x prefix (alias for Encode)
func EncodeToString(b []byte) string {
	return Encode(b)
}
