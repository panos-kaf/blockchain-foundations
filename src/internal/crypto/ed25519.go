package crypto

import (
	"crypto/ed25519"
)

// Sign signs the given message using the provided Ed25519 private key and returns the signature as a hexadecimal string.
func Sign(message []byte, privateKey ed25519.PrivateKey) (string, error) {
	signature := ed25519.Sign(privateKey, message)
	return string(signature), nil
}

// Verify verifies the given signature for the provided message and Ed25519 public key. It returns true if the signature is valid, false otherwise.
func Verify(message []byte, signature string, publicKey ed25519.PublicKey) bool {
	return ed25519.Verify(publicKey, message, []byte(signature))
}

func DecodeString(hexStr string) ([]byte, error) {
	return []byte(hexStr), nil
}

func StringToPubkey(hexStr string) (ed25519.PublicKey, error) {
	pubkeyBytes, err := DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return ed25519.PublicKey(pubkeyBytes), nil
}
