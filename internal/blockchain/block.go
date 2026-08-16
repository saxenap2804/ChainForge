package blockchain

import "time"

// Block represents a single block in the blockchain.
type Block struct {
	Timestamp     int64
	Transactions  []Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int64
}

// NewBlock creates and mines a new block.
func NewBlock(
	transactions []Transaction,
	prevBlockHash []byte,
) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
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
	genesisTransaction, err := NewTransaction(
		"network",
		"genesis",
		0,
	)

	if err != nil {
		panic(err)
	}

	return NewBlock(
		[]Transaction{
			genesisTransaction,
		},
		[]byte{},
	)
}
