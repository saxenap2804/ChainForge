package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

// Transaction represents a transfer of value
// between two addresses.
type Transaction struct {
	ID       string  `json:"id"`
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
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

// Serialize converts the transaction into bytes.
// ID is intentionally excluded so the ID itself
// does not affect transaction hashing.
func (tx Transaction) Serialize() []byte {
	copyTx := tx
	copyTx.ID = ""

	data, err := json.Marshal(copyTx)

	if err != nil {
		panic(err)
	}

	return data
}

// calculateID generates the deterministic
// SHA-256 identifier for the transaction.
func (tx Transaction) calculateID() string {
	hash := sha256.Sum256(
		tx.Serialize(),
	)

	return hex.EncodeToString(
		hash[:],
	)
}

// Hash returns the SHA-256 transaction hash.
func (tx Transaction) Hash() []byte {
	hash := sha256.Sum256(
		tx.Serialize(),
	)

	return hash[:]
}
