package blockchain

import (
	"bytes"
	"testing"
)

func TestNewGenesisBlock(t *testing.T) {
	block := NewGenesisBlock()

	if block == nil {
		t.Fatal("expected genesis block, got nil")
	}

	if string(block.Data) != "Genesis Block" {
		t.Fatalf(
			"expected genesis block data, got %q",
			string(block.Data),
		)
	}

	if len(block.Hash) == 0 {
		t.Fatal("expected genesis block hash to be set")
	}
}

func TestNewBlockchain(t *testing.T) {
	chain := NewBlockchain()

	if chain == nil {
		t.Fatal("expected blockchain, got nil")
	}

	if len(chain.Blocks) != 1 {
		t.Fatalf(
			"expected 1 genesis block, got %d",
			len(chain.Blocks),
		)
	}
}

func TestAddBlock(t *testing.T) {
	chain := NewBlockchain()

	chain.AddBlock("Alice sends 10 coins to Bob")

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
			"new block does not reference previous block hash",
		)
	}
}

func TestBlockchainIsValid(t *testing.T) {
	chain := NewBlockchain()

	chain.AddBlock("Alice sends 10 coins to Bob")
	chain.AddBlock("Bob sends 4 coins to Charlie")

	if !chain.IsValid() {
		t.Fatal("expected blockchain to be valid")
	}
}

func TestBlockchainDetectsTampering(t *testing.T) {
	chain := NewBlockchain()

	chain.AddBlock("Alice sends 10 coins to Bob")
	chain.AddBlock("Bob sends 4 coins to Charlie")

	chain.Blocks[1].Data = []byte(
		"Alice sends 1000 coins to Mallory",
	)

	if chain.IsValid() {
		t.Fatal(
			"expected tampered blockchain to be invalid",
		)
	}
}

func TestBlockHashChangesWithData(t *testing.T) {
	previousHash := []byte("previous")

	blockOne := NewBlock(
		"Transaction A",
		previousHash,
	)

	blockTwo := NewBlock(
		"Transaction B",
		previousHash,
	)

	if bytes.Equal(
		blockOne.Hash,
		blockTwo.Hash,
	) {
		t.Fatal(
			"expected different data to produce different hashes",
		)
	}
}
