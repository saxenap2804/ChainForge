package main

import (
	"fmt"

	"github.com/saxenap2804/ChainForge/internal/blockchain"
)

func printChain(
	chain *blockchain.Blockchain,
) {
	fmt.Println("ChainForge")
	fmt.Println("==========")

	for index, block := range chain.Blocks {
		fmt.Printf("\nBlock %d\n", index)

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

		fmt.Println("Transactions:")

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

			if len(tx.Signature) > 0 {
				fmt.Println(
					"  Signature: verified",
				)
			} else {
				fmt.Println(
					"  Signature: genesis",
				)
			}
		}
	}
}

func main() {
	// Create wallets.
	alice, err := blockchain.NewWallet()
	if err != nil {
		panic(err)
	}

	bob, err := blockchain.NewWallet()
	if err != nil {
		panic(err)
	}

	charlie, err := blockchain.NewWallet()
	if err != nil {
		panic(err)
	}

	fmt.Println("Wallets")
	fmt.Println("=======")
	fmt.Printf("Alice:   %s\n", alice.Address)
	fmt.Printf("Bob:     %s\n", bob.Address)
	fmt.Printf("Charlie: %s\n\n", charlie.Address)

	// Create blockchain.
	chain := blockchain.NewBlockchain()

	// Alice -> Bob.
	tx1, err := blockchain.NewTransaction(
		alice.Address,
		bob.Address,
		10,
	)

	if err != nil {
		panic(err)
	}

	if err := tx1.Sign(alice); err != nil {
		panic(err)
	}

	// Bob -> Charlie.
	tx2, err := blockchain.NewTransaction(
		bob.Address,
		charlie.Address,
		4,
	)

	if err != nil {
		panic(err)
	}

	if err := tx2.Sign(bob); err != nil {
		panic(err)
	}

	// Mine transactions into blocks.
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

	printChain(chain)

	fmt.Printf(
		"\nBlockchain valid before tampering: %t\n",
		chain.IsValid(),
	)

	// Demonstrate signature + blockchain tamper detection.
	fmt.Println("\nTampering with Alice's transaction...")

	chain.Blocks[1].
		Transactions[0].
		Amount = 1000

	fmt.Printf(
		"Blockchain valid after tampering: %t\n",
		chain.IsValid(),
	)
}
