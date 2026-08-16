package blockchain

import (
	"math/big"
	"testing"
)

func TestProofOfWorkValidBlock(t *testing.T) {
	block := NewBlock(
		"Alice sends 10 coins to Bob",
		[]byte{},
	)

	pow := NewProofOfWork(block)

	if !pow.Validate() {
		t.Fatal("expected mined block to satisfy proof of work")
	}
}

func TestProofOfWorkDetectsTamperedData(t *testing.T) {
	block := NewBlock(
		"Alice sends 10 coins to Bob",
		[]byte{},
	)

	block.Data = []byte(
		"Alice sends 1000 coins to Mallory",
	)

	pow := NewProofOfWork(block)

	if pow.Validate() {
		t.Fatal(
			"expected modified block data to invalidate proof of work",
		)
	}
}

func TestProofOfWorkDetectsTamperedNonce(t *testing.T) {
	block := NewBlock(
		"Alice sends 10 coins to Bob",
		[]byte{},
	)

	block.Nonce++

	pow := NewProofOfWork(block)

	if pow.Validate() {
		t.Fatal(
			"expected modified nonce to invalidate proof of work",
		)
	}
}

func TestProofOfWorkHashMeetsTarget(t *testing.T) {
	block := NewBlock(
		"Test transaction",
		[]byte{},
	)

	pow := NewProofOfWork(block)

	var hashInt big.Int
	hashInt.SetBytes(block.Hash)

	if hashInt.Cmp(pow.target) >= 0 {
		t.Fatal(
			"expected block hash to be below proof-of-work target",
		)
	}
}
