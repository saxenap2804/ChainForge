package blockchain

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// Block represents a single block in the blockchain.
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

// calculateHash computes the SHA-256 hash for a block.
func (b *Block) calculateHash() []byte {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))

	headers := bytes.Join(
		[][]byte{
			b.PrevBlockHash,
			b.Data,
			timestamp,
		},
		[]byte{},
	)

	hash := sha256.Sum256(headers)

	return hash[:]
}

// NewBlock creates a new block and calculates its hash.
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
	}

	block.Hash = block.calculateHash()

	return block
}

// NewGenesisBlock creates the first block in the chain.
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
