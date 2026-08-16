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

// AddBlock appends a new mined block to the blockchain.
func (bc *Blockchain) AddBlock(data string) {
	previousBlock := bc.Blocks[len(bc.Blocks)-1]

	newBlock := NewBlock(
		data,
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

		// Verify that this block references
		// the actual hash of the previous block.
		if !bytes.Equal(
			currentBlock.PrevBlockHash,
			previousBlock.Hash,
		) {
			return false
		}

		// Verify that the block still satisfies
		// the Proof-of-Work requirement.
		pow := NewProofOfWork(currentBlock)

		if !pow.Validate() {
			return false
		}
	}

	return true
}
