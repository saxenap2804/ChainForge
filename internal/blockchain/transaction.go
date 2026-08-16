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

// NewTransaction creates and validates a normal transaction.
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

// NewCoinbaseTransaction creates new coins
// and assigns them to a recipient.
//
// Coinbase transactions are issued by the network
// and do not require a sender signature.
func NewCoinbaseTransaction(
	receiver string,
	amount float64,
) (Transaction, error) {
	receiver = strings.TrimSpace(receiver)

	if receiver == "" {
		return Transaction{},
			errors.New("receiver cannot be empty")
	}

	if amount <= 0 {
		return Transaction{},
			errors.New("coinbase amount must be positive")
	}

	tx := Transaction{
		Sender:   "network",
		Receiver: receiver,
		Amount:   amount,
	}

	tx.ID = tx.calculateID()

	return tx, nil
}

// Validate checks the basic transaction rules.
func (tx Transaction) Validate() error {
	if tx.Sender == "" {
		return errors.New("sender cannot be empty")
	}

	if tx.Receiver == "" {
		return errors.New("receiver cannot be empty")
	}

	if tx.Sender == tx.Receiver {
		return errors.New(
			"sender and receiver cannot be the same",
		)
	}

	if tx.Amount < 0 {
		return errors.New(
			"amount cannot be negative",
		)
	}

	return nil
}

// IsCoinbase reports whether this transaction
// was issued by the network.
func (tx Transaction) IsCoinbase() bool {
	return tx.Sender == "network"
}

// Serialize returns the canonical transaction payload.
//
// ID, public key, and signature are intentionally
// excluded so the payload remains stable for hashing
// and signing.
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

// calculateID generates a deterministic SHA-256
// identifier from the transaction contents.
func (tx Transaction) calculateID() string {
	hash := sha256.Sum256(
		tx.Serialize(),
	)

	return hex.EncodeToString(
		hash[:],
	)
}

// Hash returns the SHA-256 hash of the canonical
// transaction payload.
func (tx Transaction) Hash() []byte {
	hash := sha256.Sum256(
		tx.Serialize(),
	)

	return hash[:]
}

// Sign signs a normal transaction with the
// sender wallet's ECDSA private key.
func (tx *Transaction) Sign(
	wallet *Wallet,
) error {
	if tx.IsCoinbase() {
		return errors.New(
			"coinbase transactions do not require signatures",
		)
	}

	if wallet == nil {
		return errors.New(
			"wallet cannot be nil",
		)
	}

	if wallet.PrivateKey == nil {
		return errors.New(
			"wallet private key cannot be nil",
		)
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

	// Encode R and S using fixed 32-byte fields.
	// This avoids ambiguity when either integer
	// contains leading zero bytes.
	signature := make(
		[]byte,
		64,
	)

	rBytes := r.Bytes()
	sBytes := s.Bytes()

	copy(
		signature[32-len(rBytes):32],
		rBytes,
	)

	copy(
		signature[64-len(sBytes):],
		sBytes,
	)

	tx.PublicKey = append(
		[]byte(nil),
		wallet.PublicKey...,
	)

	tx.Signature = signature

	// The ID is based on the unsigned payload,
	// so this remains deterministic.
	tx.ID = tx.calculateID()

	return nil
}

// Verify checks the ECDSA signature of
// a normal transaction.
func (tx Transaction) Verify() bool {
	// Coinbase transactions are handled separately
	// by blockchain-level validation.
	if tx.IsCoinbase() {
		return false
	}

	if len(tx.PublicKey) == 0 {
		return false
	}

	if len(tx.Signature) != 64 {
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

	r := new(big.Int).SetBytes(
		tx.Signature[:32],
	)

	s := new(big.Int).SetBytes(
		tx.Signature[32:],
	)

	return ecdsa.Verify(
		&publicKey,
		tx.Hash(),
		r,
		s,
	)
}
