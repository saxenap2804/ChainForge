package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"strings"
)

// Transaction represents a transfer of value
// between two addresses.
type Transaction struct {
	ID        string  `json:"id"`
	Sender    string  `json:"sender"`
	Receiver  string  `json:"receiver"`
	Amount    float64 `json:"amount"`
	PublicKey []byte  `json:"public_key,omitempty"`
	Signature []byte  `json:"signature,omitempty"`
}

// NewTransaction creates and validates a transaction.
func NewTransaction(
	sender string,
	receiver string,
	amount float64,
) (Transaction, error) {
	tx := Transaction{
		Sender:   strings.TrimSpace(sender),
		Receiver: strings.TrimSpace(receiver),
		Amount:   amount,
	}

	if err := tx.Validate(); err != nil {
		return Transaction{}, err
	}

	tx.ID = tx.calculateID()

	return tx, nil
}

// Validate checks basic transaction rules.
func (tx Transaction) Validate() error {
	if tx.Sender == "" {
		return errors.New("sender cannot be empty")
	}

	if tx.Receiver == "" {
		return errors.New("receiver cannot be empty")
	}

	if tx.Sender == tx.Receiver {
		return errors.New("sender and receiver cannot be the same")
	}

	if tx.Amount < 0 {
		return errors.New("amount cannot be negative")
	}

	return nil
}

// Serialize returns the canonical transaction payload.
// ID, signature, and public key are excluded so signing
// remains deterministic.
func (tx Transaction) Serialize() []byte {
	copyTx := tx

	copyTx.ID = ""
	copyTx.PublicKey = nil
	copyTx.Signature = nil

	data, err := json.Marshal(copyTx)

	if err != nil {
		panic(err)
	}

	return data
}

// calculateID generates the deterministic transaction ID.
func (tx Transaction) calculateID() string {
	hash := sha256.Sum256(
		tx.Serialize(),
	)

	return hex.EncodeToString(
		hash[:],
	)
}

// Hash returns the transaction payload hash.
func (tx Transaction) Hash() []byte {
	hash := sha256.Sum256(
		tx.Serialize(),
	)

	return hash[:]
}

// Sign signs the transaction using the sender's wallet.
func (tx *Transaction) Sign(wallet *Wallet) error {
	if wallet == nil || wallet.PrivateKey == nil {
		return errors.New("invalid wallet")
	}

	if tx.Sender != wallet.Address {
		return errors.New(
			"wallet address does not match transaction sender",
		)
	}

	hash := tx.Hash()

	r, s, err := ecdsa.Sign(
		rand.Reader,
		wallet.PrivateKey,
		hash,
	)

	if err != nil {
		return err
	}

	signature := append(
		r.Bytes(),
		s.Bytes()...,
	)

	tx.PublicKey = wallet.PublicKey
	tx.Signature = signature
	tx.ID = tx.calculateID()

	return nil
}

// Verify checks the ECDSA signature of a transaction.
func (tx Transaction) Verify() bool {
	if len(tx.PublicKey) == 0 {
		return false
	}

	if len(tx.Signature) == 0 {
		return false
	}

	x, y := elliptic.Unmarshal(
		elliptic.P256(),
		tx.PublicKey,
	)

	if x == nil || y == nil {
		return false
	}

	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	if len(tx.Signature)%2 != 0 {
		return false
	}

	half := len(tx.Signature) / 2

	r := new(big.Int).SetBytes(
		tx.Signature[:half],
	)

	s := new(big.Int).SetBytes(
		tx.Signature[half:],
	)

	return ecdsa.Verify(
		&publicKey,
		tx.Hash(),
		r,
		s,
	)
}
