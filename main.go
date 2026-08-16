package main

import (
	"fmt"
	"log"

	"github.com/saxenap2804/ChainForge/internal/blockchain"
)

func printChain(chain *blockchain.Blockchain) {
	fmt.Println("\nChainForge")
	fmt.Println("==========")

	for index, block := range chain.Blocks {
		fmt.Printf("\nBlock %d\n", index)
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Previous Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Hash: %x\n", block.Hash)

		fmt.Println("Transactions:")

		for _, tx := range block.Transactions {
			fmt.Printf(
				"  %s -> %s : %.2f coins\n",
				tx.Sender,
				tx.Receiver,
				tx.Amount,
			)
		}
	}
}

func main() {
	// Open an existing blockchain or create a new one.
	chain, err := blockchain.OpenBlockchain(
		"chainforge.db",
	)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := chain.Close(); err != nil {
			log.Printf(
				"failed to close blockchain: %v",
				err,
			)
		}
	}()

	fmt.Println("ChainForge Persistent Node")
	fmt.Println("==========================")

	fmt.Printf(
		"Loaded %d blocks from disk.\n",
		len(chain.Blocks),
	)

	fmt.Printf(
		"Blockchain valid: %t\n",
		chain.IsValid(),
	)

	printChain(chain)
}
