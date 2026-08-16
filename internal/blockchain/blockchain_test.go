package blockchain

import (
	"bytes"
	"testing"
)

func mustTransaction(
	t *testing.T,
	sender string,
	receiver string,
	amount float64,
) Transaction {
	t.Helper()

	tx, err := NewTransaction(
		sender,
		receiver,
		amount,
	)

	if err != nil {
		t.Fatalf(
			"unexpected transaction error: %v",
			err,
		)
	}

	return tx
}

func TestNewGenesisBlock(
	t *testing.T,
) {
	block := NewGenesisBlock()

	if block == nil {
		t.Fatal(
			"expected genesis block, got nil",
		)
	}

	if len(block.Transactions) != 1 {
		t.Fatalf(
			"expected 1 genesis transaction, got %d",
			len(block.Transactions),
		)
	}

	tx := block.Transactions[0]

	if tx.Sender != "network" {
		t.Fatalf(
			"expected genesis sender network, got %q",
			tx.Sender,
		)
	}

	if tx.Receiver != "genesis" {
		t.Fatalf(
			"expected genesis receiver genesis, got %q",
			tx.Receiver,
		)
	}

	if tx.Amount != 0 {
		t.Fatalf(
			"expected genesis amount 0, got %f",
			tx.Amount,
		)
	}

	if tx.ID == "" {
		t.Fatal(
			"expected genesis transaction ID",
		)
	}

	if len(block.Hash) == 0 {
		t.Fatal(
			"expected genesis block hash",
		)
	}
}

func TestNewBlockchain(
	t *testing.T,
) {
	chain := NewBlockchain()

	if chain == nil {
		t.Fatal(
			"expected blockchain, got nil",
		)
	}

	if len(chain.Blocks) != 1 {
		t.Fatalf(
			"expected 1 block, got %d",
			len(chain.Blocks),
		)
	}
}

func TestAddBlock(
	t *testing.T,
) {
	chain := NewBlockchain()

	tx := mustTransaction(
		t,
		"Alice",
		"Bob",
		10,
	)

	chain.AddBlock(
		[]Transaction{
			tx,
		},
	)

	if len(chain.Blocks) != 2 {
		t.Fatalf(
			"expected 2 blocks, got %d",
			len(chain.Blocks),
		)
	}

	previous := chain.Blocks[0]
	current := chain.Blocks[1]

	if !bytes.Equal(
		current.PrevBlockHash,
		previous.Hash,
	) {
		t.Fatal(
			"block does not reference previous block",
		)
	}
}

func TestBlockchainIsValid(
	t *testing.T,
) {
	chain := NewBlockchain()

	tx1 := mustTransaction(
		t,
		"Alice",
		"Bob",
		10,
	)

	tx2 := mustTransaction(
		t,
		"Bob",
		"Charlie",
		4,
	)

	chain.AddBlock(
		[]Transaction{
			tx1,
		},
	)

	chain.AddBlock(
		[]Transaction{
			tx2,
		},
	)

	if !chain.IsValid() {
		t.Fatal(
			"expected blockchain to be valid",
		)
	}
}

func TestBlockchainDetectsTampering(
	t *testing.T,
) {
	chain := NewBlockchain()

	tx := mustTransaction(
		t,
		"Alice",
		"Bob",
		10,
	)

	chain.AddBlock(
		[]Transaction{
			tx,
		},
	)

	chain.Blocks[1].
		Transactions[0].
		Amount = 1000

	if chain.IsValid() {
		t.Fatal(
			"expected tampered blockchain to be invalid",
		)
	}
}

func TestBlockHashChangesWithTransactions(
	t *testing.T,
) {
	previousHash := []byte(
		"previous",
	)

	tx1 := mustTransaction(
		t,
		"Alice",
		"Bob",
		10,
	)

	tx2 := mustTransaction(
		t,
		"Alice",
		"Charlie",
		10,
	)

	blockOne := NewBlock(
		[]Transaction{
			tx1,
		},
		previousHash,
	)

	blockTwo := NewBlock(
		[]Transaction{
			tx2,
		},
		previousHash,
	)

	if bytes.Equal(
		blockOne.Hash,
		blockTwo.Hash,
	) {
		t.Fatal(
			"different transactions should produce different hashes",
		)
	}
}

func TestBlockCanContainMultipleTransactions(
	t *testing.T,
) {
	tx1 := mustTransaction(
		t,
		"Alice",
		"Bob",
		10,
	)

	tx2 := mustTransaction(
		t,
		"Bob",
		"Charlie",
		4,
	)

	block := NewBlock(
		[]Transaction{
			tx1,
			tx2,
		},
		[]byte{},
	)

	if len(block.Transactions) != 2 {
		t.Fatalf(
			"expected 2 transactions, got %d",
			len(block.Transactions),
		)
	}
}
