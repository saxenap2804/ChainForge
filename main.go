package main

import (
	"fmt"

	"github.com/saxenap2804/ChainForge/internal/blockchain"
)

// printChain displays the complete blockchain.
func printChain(
	chain *blockchain.Blockchain,
) {
	fmt.Println("ChainForge")
	fmt.Println("==========")

	for index, block := range chain.Blocks {
		fmt.Printf(
			"\nBlock %d\n",
			index,
		)

		fmt.Printf(
			"Timestamp: %d\n",
			block.Timestamp,
		)

		fmt.Printf(
			"Previous Hash: %x\n",
			block.PrevBlockHash,
		)

		fmt.Printf(
			"Nonce: %d\n",
			block.Nonce,
		)

		fmt.Printf(
			"Hash: %x\n",
			block.Hash,
		)

		fmt.Println(
			"Transactions:",
		)

		for _, tx := range block.Transactions {
			fmt.Printf(
				"  ID: %s\n",
				tx.ID,
			)

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
	chain := blockchain.NewBlockchain()

	tx1, err := blockchain.NewTransaction(
		"Alice",
		"Bob",
		10,
	)

	if err != nil {
		panic(err)
	}

	tx2, err := blockchain.NewTransaction(
		"Bob",
		"Charlie",
		4,
	)

	if err != nil {
		panic(err)
	}

	chain.AddBlock(
		[]blockchain.Transaction{
			tx1,
		},
	)

	chain.AddBlock(
		[]blockchain.Transaction{
			tx2,
		},
	)

	printChain(
		chain,
	)

	fmt.Printf(
		"\nBlockchain valid before tampering: %t\n",
		chain.IsValid(),
	)

	fmt.Println(
		"\nTampering with Block 1...",
	)

	chain.Blocks[1].
		Transactions[0].
		Amount = 1000

	fmt.Printf(
		"Blockchain valid after tampering: %t\n",
		chain.IsValid(),
	)
}
