package blockchain

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestStorageSaveAndLoadBlock(
	t *testing.T,
) {
	tempDir := t.TempDir()

	dbPath := filepath.Join(
		tempDir,
		"chainforge.db",
	)

	storage, err := OpenStorage(
		dbPath,
	)

	if err != nil {
		t.Fatal(err)
	}

	defer storage.Close()

	tx, err := NewCoinbaseTransaction(
		"test-wallet",
		50,
	)

	if err != nil {
		t.Fatal(err)
	}

	block := NewBlock(
		[]Transaction{
			tx,
		},
		[]byte{},
	)

	if err := storage.SaveBlock(
		block,
	); err != nil {
		t.Fatal(err)
	}

	loaded, err := storage.LoadBlock(
		block.Hash,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(
		loaded.Hash,
		block.Hash,
	) {
		t.Fatal(
			"loaded block hash does not match original",
		)
	}

	lastHash, err := storage.LastHash()

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(
		lastHash,
		block.Hash,
	) {
		t.Fatal(
			"stored latest hash does not match block hash",
		)
	}
}
