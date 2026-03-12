package crypto

import (
	"golang.org/x/crypto/blake2s"
)

// Hash computes the BLAKE2s hash of the input data and returns it as a hexadecimal string.
func Hash(data []byte) (string, error) {
	hasher, err := blake2s.New256(nil)
	if err != nil {
		return "", err
	}
	hasher.Write(data)
	return string(hasher.Sum(nil)), nil
}

// HashString is a convenience function that takes a string input, computes its BLAKE2s hash, and returns the hash as a hexadecimal string.
func HashString(s string) (string, error) {
	return Hash([]byte(s))
}

// HashBytes is a convenience function that takes a byte slice input, computes its BLAKE2s hash, and returns the hash as a hexadecimal string.
func HashBytes(b []byte) (string, error) {
	return Hash(b)
}
