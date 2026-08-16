package blockchain

import (
	"math/big"
	"testing"
)

func TestProofOfWorkValidBlock(
	t *testing.T,
) {
	tx := mustTransaction(
		t,
		"Alice",
		"Bob",
		10,
	)

	block := NewBlock(
		[]Transaction{
			tx,
		},
		[]byte{},
	)

	pow := NewProofOfWork(
		block,
	)

	if !pow.Validate() {
		t.Fatal(
			"expected mined block to satisfy proof of work",
		)
	}
}

func TestProofOfWorkDetectsTamperedTransaction(
	t *testing.T,
) {
	tx := mustTransaction(
		t,
		"Alice",
		"Bob",
		10,
	)

	block := NewBlock(
		[]Transaction{
			tx,
		},
		[]byte{},
	)

	block.Transactions[0].Amount = 1000

	pow := NewProofOfWork(
		block,
	)

	if pow.Validate() {
		t.Fatal(
			"expected modified transaction to invalidate proof of work",
		)
	}
}

func TestProofOfWorkDetectsTamperedNonce(
	t *testing.T,
) {
	tx := mustTransaction(
		t,
		"Alice",
		"Bob",
		10,
	)

	block := NewBlock(
		[]Transaction{
			tx,
		},
		[]byte{},
	)

	block.Nonce++

	pow := NewProofOfWork(
		block,
	)

	if pow.Validate() {
		t.Fatal(
			"expected modified nonce to invalidate proof of work",
		)
	}
}

func TestProofOfWorkHashMeetsTarget(
	t *testing.T,
) {
	tx := mustTransaction(
		t,
		"Tester",
		"Receiver",
		1,
	)

	block := NewBlock(
		[]Transaction{
			tx,
		},
		[]byte{},
	)

	pow := NewProofOfWork(
		block,
	)

	var hashInt big.Int

	hashInt.SetBytes(
		block.Hash,
	)

	if hashInt.Cmp(
		pow.target,
	) >= 0 {
		t.Fatal(
			"expected block hash below proof-of-work target",
		)
	}
}
