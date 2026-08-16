package blockchain

import "time"

// Block represents a single block in the blockchain.
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int64
}

// NewBlock creates a new block and mines it
// using the Proof-of-Work algorithm.
func NewBlock(
	data string,
	prevBlockHash []byte,
) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}

	pow := NewProofOfWork(block)

	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates the first block
// in the blockchain.
func NewGenesisBlock() *Block {
	return NewBlock(
		"Genesis Block",
		[]byte{},
	)
}
