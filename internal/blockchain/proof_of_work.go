package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"
)

const targetBits = 18

// ProofOfWork represents the mining state for a block.
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork creates a proof-of-work instance
// for the supplied block.
func NewProofOfWork(
	block *Block,
) *ProofOfWork {
	target := big.NewInt(1)

	target.Lsh(
		target,
		uint(256-targetBits),
	)

	return &ProofOfWork{
		block:  block,
		target: target,
	}
}

// intToBytes converts an int64 to bytes.
func intToBytes(
	num int64,
) []byte {
	buffer := new(
		bytes.Buffer,
	)

	_ = binary.Write(
		buffer,
		binary.BigEndian,
		num,
	)

	return buffer.Bytes()
}

// hashTransactions combines all transaction hashes
// into a single SHA-256 digest.
func (pow *ProofOfWork) hashTransactions() []byte {
	var transactionHashes []byte

	for _, tx := range pow.block.Transactions {
		transactionHashes = append(
			transactionHashes,
			tx.Hash()...,
		)
	}

	hash := sha256.Sum256(
		transactionHashes,
	)

	return hash[:]
}

// prepareData creates the byte sequence hashed
// during Proof-of-Work mining.
func (pow *ProofOfWork) prepareData(
	nonce int64,
) []byte {
	return bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.hashTransactions(),
			intToBytes(
				pow.block.Timestamp,
			),
			intToBytes(
				int64(targetBits),
			),
			intToBytes(
				nonce,
			),
		},
		[]byte{},
	)
}

// Run searches for a nonce whose block hash
// satisfies the current difficulty target.
func (pow *ProofOfWork) Run() (
	int64,
	[]byte,
) {
	var hashInt big.Int
	var hash [32]byte

	var nonce int64

	for nonce < math.MaxInt64 {
		data := pow.prepareData(
			nonce,
		)

		hash = sha256.Sum256(
			data,
		)

		hashInt.SetBytes(
			hash[:],
		)

		if hashInt.Cmp(
			pow.target,
		) == -1 {
			break
		}

		nonce++
	}

	return nonce, hash[:]
}

// Validate verifies that a block's stored nonce
// still satisfies the Proof-of-Work requirement.
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(
		pow.block.Nonce,
	)

	hash := sha256.Sum256(
		data,
	)

	hashInt.SetBytes(
		hash[:],
	)

	return hashInt.Cmp(
		pow.target,
	) == -1
}
