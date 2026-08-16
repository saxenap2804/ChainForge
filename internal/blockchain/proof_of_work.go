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
// for the given block.
func NewProofOfWork(block *Block) *ProofOfWork {
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

// intToBytes converts an int64 into bytes.
func intToBytes(num int64) []byte {
	buffer := new(bytes.Buffer)

	_ = binary.Write(
		buffer,
		binary.BigEndian,
		num,
	)

	return buffer.Bytes()
}

// prepareData creates the byte sequence that is hashed
// during mining.
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	return bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			intToBytes(pow.block.Timestamp),
			intToBytes(int64(targetBits)),
			intToBytes(nonce),
		},
		[]byte{},
	)
}

// Run searches for a nonce whose hash satisfies
// the current difficulty target.
func (pow *ProofOfWork) Run() (int64, []byte) {
	var hashInt big.Int
	var hash [32]byte

	var nonce int64

	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)

		hash = sha256.Sum256(data)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		}

		nonce++
	}

	return nonce, hash[:]
}

// Validate verifies that the block's stored nonce
// produces a hash below the required target.
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(
		pow.block.Nonce,
	)

	hash := sha256.Sum256(data)

	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}
