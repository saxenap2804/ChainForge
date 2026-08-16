package blockchain

import (
	"bytes"
	"errors"
)

// Blockchain stores an ordered collection of blocks
// and optionally persists them using Storage.
type Blockchain struct {
	Blocks  []*Block
	Storage *Storage
}

// NewBlockchain creates an in-memory blockchain
// containing the genesis block.
func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []*Block{
			NewGenesisBlock(),
		},
	}
}

// OpenBlockchain opens or creates a persistent blockchain.
func OpenBlockchain(path string) (*Blockchain, error) {
	storage, err := OpenStorage(path)
	if err != nil {
		return nil, err
	}

	chain := &Blockchain{
		Storage: storage,
	}

	lastHash, err := storage.LastHash()

	// No existing chain: create and persist genesis block.
	if err != nil {
		genesis := NewGenesisBlock()

		if err := storage.SaveBlock(genesis); err != nil {
			storage.Close()
			return nil, err
		}

		chain.Blocks = []*Block{
			genesis,
		}

		return chain, nil
	}

	// Existing chain: reconstruct it from disk.
	var reversed []*Block
	currentHash := lastHash

	for len(currentHash) > 0 {
		block, err := storage.LoadBlock(currentHash)
		if err != nil {
			storage.Close()
			return nil, err
		}

		reversed = append(
			reversed,
			block,
		)

		currentHash = block.PrevBlockHash
	}

	// Reverse blocks into genesis -> latest order.
	for i := len(reversed) - 1; i >= 0; i-- {
		chain.Blocks = append(
			chain.Blocks,
			reversed[i],
		)
	}

	if len(chain.Blocks) == 0 {
		storage.Close()
		return nil, errors.New(
			"persistent blockchain contains no blocks",
		)
	}

	return chain, nil
}

// Close closes persistent storage if enabled.
func (bc *Blockchain) Close() error {
	if bc.Storage == nil {
		return nil
	}

	return bc.Storage.Close()
}

// AddBlock validates transactions, mines a block,
// appends it to memory, and persists it when storage
// is enabled.
func (bc *Blockchain) AddBlock(
	transactions []Transaction,
) error {
	if len(transactions) == 0 {
		return errors.New(
			"block must contain at least one transaction",
		)
	}

	for _, tx := range transactions {
		if err := bc.ValidateTransaction(tx); err != nil {
			return err
		}
	}

	previousBlock := bc.Blocks[len(bc.Blocks)-1]

	newBlock := NewBlock(
		transactions,
		previousBlock.Hash,
	)

	// Persist before updating in-memory state.
	if bc.Storage != nil {
		if err := bc.Storage.SaveBlock(
			newBlock,
		); err != nil {
			return err
		}
	}

	bc.Blocks = append(
		bc.Blocks,
		newBlock,
	)

	return nil
}

// IsValid verifies chain linkage, transaction
// integrity, signatures, and Proof of Work.
func (bc *Blockchain) IsValid() bool {
	if len(bc.Blocks) == 0 {
		return false
	}

	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		previousBlock := bc.Blocks[i-1]

		if !bytes.Equal(
			currentBlock.PrevBlockHash,
			previousBlock.Hash,
		) {
			return false
		}

		for _, tx := range currentBlock.Transactions {
			if err := tx.Validate(); err != nil {
				return false
			}

			if tx.ID != tx.calculateID() {
				return false
			}

			if !tx.IsCoinbase() && !tx.Verify() {
				return false
			}
		}

		pow := NewProofOfWork(
			currentBlock,
		)

		if !pow.Validate() {
			return false
		}
	}

	return true
}
