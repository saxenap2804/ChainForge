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

	chain := blockchain.NewBlockchain()

	// Give Alice starting funds.
	reward, err := blockchain.NewCoinbaseTransaction(
		alice.Address,
		50,
	)

	if err != nil {
		panic(err)
	}

	if err := chain.AddBlock(
		[]blockchain.Transaction{
			reward,
		},
	); err != nil {
		panic(err)
	}

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

	if err := chain.AddBlock(
		[]blockchain.Transaction{
			tx1,
		},
	); err != nil {
		panic(err)
	}

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

	if err := chain.AddBlock(
		[]blockchain.Transaction{
			tx2,
		},
	); err != nil {
		panic(err)
	}

	printChain(chain)

	fmt.Println("\nBalances")
	fmt.Println("========")

	fmt.Printf(
		"Alice: %.2f\n",
		chain.BalanceOf(alice.Address),
	)

	fmt.Printf(
		"Bob: %.2f\n",
		chain.BalanceOf(bob.Address),
	)

	fmt.Printf(
		"Charlie: %.2f\n",
		chain.BalanceOf(charlie.Address),
	)

	fmt.Printf(
		"\nBlockchain valid: %t\n",
		chain.IsValid(),
	)

	// Try to overspend.
	fmt.Println(
		"\nAttempting invalid overspend from Bob...",
	)

	badTx, err := blockchain.NewTransaction(
		bob.Address,
		charlie.Address,
		100,
	)

	if err != nil {
		panic(err)
	}

	if err := badTx.Sign(bob); err != nil {
		panic(err)
	}

	err = chain.AddBlock(
		[]blockchain.Transaction{
			badTx,
		},
	)

	if err != nil {
		fmt.Printf(
			"Overspend rejected: %v\n",
			err,
		)
	}
}
