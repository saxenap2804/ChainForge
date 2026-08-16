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

// AddBlock validates transactions and appends
// a newly mined block to the chain.
func (bc *Blockchain) AddBlock(
	transactions []Transaction,
) error {
	if len(transactions) == 0 {
		return nil
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
