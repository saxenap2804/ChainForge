package blockchain

import "time"

const genesisTimestamp int64 = 1700000000

// Block represents a single block in the blockchain.
type Block struct {
	Timestamp     int64
	Transactions  []Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int64
}

// NewBlock creates and mines a normal block.
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

// NewGenesisBlock creates the deterministic first
// block shared by every ChainForge node.
func NewGenesisBlock() *Block {
	genesisTransaction, err := NewTransaction(
		"network",
		"genesis",
		0,
	)

	if err != nil {
		panic(err)
	}

	block := &Block{
		Timestamp: genesisTimestamp,
		Transactions: []Transaction{
			genesisTransaction,
		},
		PrevBlockHash: []byte{},
		Hash:          []byte{},
		Nonce:         0,
	}

	pow := NewProofOfWork(block)

	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}
