package blockchain

import (
	"path/filepath"
	"testing"
)

func TestPersistentBlockchainReopens(
	t *testing.T,
) {
	tempDir := t.TempDir()

	dbPath := filepath.Join(
		tempDir,
		"chainforge.db",
	)

	chain, err := OpenBlockchain(
		dbPath,
	)

	if err != nil {
		t.Fatal(err)
	}

	wallet, err := NewWallet()
	if err != nil {
		t.Fatal(err)
	}

	reward, err := NewCoinbaseTransaction(
		wallet.Address,
		50,
	)

	if err != nil {
		t.Fatal(err)
	}

	if err := chain.AddBlock(
		[]Transaction{
			reward,
		},
	); err != nil {
		t.Fatal(err)
	}

	expectedBlockCount := len(
		chain.Blocks,
	)

	if err := chain.Close(); err != nil {
		t.Fatal(err)
	}

	reopened, err := OpenBlockchain(
		dbPath,
	)

	if err != nil {
		t.Fatal(err)
	}

	defer reopened.Close()

	if len(reopened.Blocks) != expectedBlockCount {
		t.Fatalf(
			"expected %d blocks after reopen, got %d",
			expectedBlockCount,
			len(reopened.Blocks),
		)
	}

	if !reopened.IsValid() {
		t.Fatal(
			"expected reopened blockchain to be valid",
		)
	}
}
