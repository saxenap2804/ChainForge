package main

import (
	"fmt"

	"github.com/saxenap2804/ChainForge/internal/blockchain"
)

func printChain(chain *blockchain.Blockchain) {
	fmt.Println("ChainForge")
	fmt.Println("==========")

	for index, block := range chain.Blocks {
		fmt.Printf("\nBlock %d\n", index)
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Previous Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
	}
}

func main() {
	chain := blockchain.NewBlockchain()

	chain.AddBlock(
		"Alice sends 10 coins to Bob",
	)

	chain.AddBlock(
		"Bob sends 4 coins to Charlie",
	)

	printChain(chain)

	fmt.Printf(
		"\nBlockchain valid before tampering: %t\n",
		chain.IsValid(),
	)

	// Simulate malicious modification of an existing block.
	fmt.Println("\nTampering with Block 1...")

	chain.Blocks[1].Data = []byte(
		"Alice sends 1000 coins to Mallory",
	)

	fmt.Printf(
		"Blockchain valid after tampering: %t\n",
		chain.IsValid(),
	)
}
