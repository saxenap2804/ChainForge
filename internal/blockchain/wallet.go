package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// Wallet contains an ECDSA key pair and derived address.
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
	Address    string
}

// NewWallet generates a new ECDSA wallet.
func NewWallet() (*Wallet, error) {
	privateKey, err := ecdsa.GenerateKey(
		elliptic.P256(),
		rand.Reader,
	)

	if err != nil {
		return nil, err
	}

	publicKey := elliptic.Marshal(
		elliptic.P256(),
		privateKey.PublicKey.X,
		privateKey.PublicKey.Y,
	)

	address := generateAddress(publicKey)

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}, nil
}

// generateAddress creates a simple SHA-256-based wallet address.
func generateAddress(publicKey []byte) string {
	hash := sha256.Sum256(publicKey)

	return hex.EncodeToString(
		hash[:20],
	)
}
