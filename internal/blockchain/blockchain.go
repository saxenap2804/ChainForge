package blockchain

import "bytes"

// Blockchain stores an ordered collection of blocks.
type Blockchain struct {
	Blocks []*Block
}

// NewBlockchain creates a blockchain containing
// the genesis block.
func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []*Block{
			NewGenesisBlock(),
		},
	}
}

// AddBlock appends a new mined block containing
// one or more transactions.
func (bc *Blockchain) AddBlock(transactions []Transaction) {
	previousBlock := bc.Blocks[len(bc.Blocks)-1]

	newBlock := NewBlock(
		transactions,
		previousBlock.Hash,
	)

	bc.Blocks = append(
		bc.Blocks,
		newBlock,
	)
}

// IsValid verifies the integrity of the blockchain.
func (bc *Blockchain) IsValid() bool {
	if len(bc.Blocks) == 0 {
		return false
	}

	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		previousBlock := bc.Blocks[i-1]

		// Verify chain linkage.
		if !bytes.Equal(
			currentBlock.PrevBlockHash,
			previousBlock.Hash,
		) {
			return false
		}

		// Validate every transaction.
		for _, tx := range currentBlock.Transactions {
			if err := tx.Validate(); err != nil {
				return false
			}

			if tx.ID != tx.calculateID() {
				return false
			}
		}

		// Validate Proof of Work.
		pow := NewProofOfWork(currentBlock)

		if !pow.Validate() {
			return false
		}
	}

	return true
}
