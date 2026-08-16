package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/saxenap2804/ChainForge/internal/blockchain"
)

const (
	databasePath = "chainforge.db"
	walletPath   = "wallets.json"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "createwallet":
		createWallet()

	case "balance":
		getBalance(os.Args[2:])

	case "fund":
		fundAddress(os.Args[2:])

	case "send":
		sendCoins(os.Args[2:])

	case "printchain":
		printBlockchain()

	default:
		fmt.Printf(
			"Unknown command: %s\n\n",
			os.Args[1],
		)

		printUsage()
	}
}

func printUsage() {
	fmt.Println("ChainForge CLI")
	fmt.Println("==============")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  createwallet")
	fmt.Println("  balance --address ADDRESS")
	fmt.Println("  fund --address ADDRESS --amount AMOUNT")
	fmt.Println("  send --from ADDRESS --to ADDRESS --amount AMOUNT")
	fmt.Println("  printchain")
}

// createWallet generates a new ECDSA wallet
// and persists it to wallets.json.
func createWallet() {
	store := blockchain.NewWalletStore(
		walletPath,
	)

	wallets, err := store.Load()
	if err != nil {
		log.Fatal(err)
	}

	wallet, err := blockchain.NewWallet()
	if err != nil {
		log.Fatal(err)
	}

	wallets[wallet.Address] = wallet

	if err := store.Save(wallets); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Wallet created successfully.")
	fmt.Printf(
		"Address: %s\n",
		wallet.Address,
	)
}

// getBalance prints the current balance
// of an address.
func getBalance(args []string) {
	command := flag.NewFlagSet(
		"balance",
		flag.ExitOnError,
	)

	address := command.String(
		"address",
		"",
		"wallet address",
	)

	if err := command.Parse(args); err != nil {
		log.Fatal(err)
	}

	if *address == "" {
		log.Fatal(
			"--address is required",
		)
	}

	chain := openChain()
	defer closeChain(chain)

	balance := chain.BalanceOf(
		*address,
	)

	fmt.Printf(
		"Balance of %s: %.2f coins\n",
		*address,
		balance,
	)
}

// fundAddress creates a coinbase transaction
// that assigns coins to an address.
func fundAddress(args []string) {
	command := flag.NewFlagSet(
		"fund",
		flag.ExitOnError,
	)

	address := command.String(
		"address",
		"",
		"wallet address",
	)

	amount := command.Float64(
		"amount",
		0,
		"amount to fund",
	)

	if err := command.Parse(args); err != nil {
		log.Fatal(err)
	}

	if *address == "" {
		log.Fatal(
			"--address is required",
		)
	}

	if *amount <= 0 {
		log.Fatal(
			"--amount must be greater than zero",
		)
	}

	chain := openChain()
	defer closeChain(chain)

	tx, err := blockchain.NewCoinbaseTransaction(
		*address,
		*amount,
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := chain.AddBlock(
		[]blockchain.Transaction{
			tx,
		},
	); err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"Funded %s with %.2f coins.\n",
		*address,
		*amount,
	)

	fmt.Printf(
		"New balance: %.2f coins\n",
		chain.BalanceOf(*address),
	)
}

// sendCoins creates, signs, validates,
// mines, and persists a transaction.
func sendCoins(args []string) {
	command := flag.NewFlagSet(
		"send",
		flag.ExitOnError,
	)

	from := command.String(
		"from",
		"",
		"sender wallet address",
	)

	to := command.String(
		"to",
		"",
		"receiver wallet address",
	)

	amount := command.Float64(
		"amount",
		0,
		"amount to send",
	)

	if err := command.Parse(args); err != nil {
		log.Fatal(err)
	}

	if *from == "" {
		log.Fatal(
			"--from is required",
		)
	}

	if *to == "" {
		log.Fatal(
			"--to is required",
		)
	}

	if *amount <= 0 {
		log.Fatal(
			"--amount must be greater than zero",
		)
	}

	store := blockchain.NewWalletStore(
		walletPath,
	)

	wallets, err := store.Load()
	if err != nil {
		log.Fatal(err)
	}

	senderWallet, exists := wallets[*from]

	if !exists {
		log.Fatal(
			"sender wallet not found",
		)
	}

	chain := openChain()
	defer closeChain(chain)

	tx, err := blockchain.NewTransaction(
		*from,
		*to,
		*amount,
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := tx.Sign(
		senderWallet,
	); err != nil {
		log.Fatal(err)
	}

	if err := chain.AddBlock(
		[]blockchain.Transaction{
			tx,
		},
	); err != nil {
		log.Fatal(err)
	}

	fmt.Println(
		"Transaction added successfully.",
	)

	fmt.Printf(
		"%s -> %s : %.2f coins\n",
		*from,
		*to,
		*amount,
	)

	fmt.Printf(
		"Sender balance: %.2f coins\n",
		chain.BalanceOf(*from),
	)

	fmt.Printf(
		"Receiver balance: %.2f coins\n",
		chain.BalanceOf(*to),
	)
}

// printBlockchain displays every block
// currently stored in the chain.
func printBlockchain() {
	chain := openChain()
	defer closeChain(chain)

	fmt.Println("ChainForge")
	fmt.Println("==========")

	fmt.Printf(
		"Blocks: %d\n",
		len(chain.Blocks),
	)

	fmt.Printf(
		"Valid: %t\n",
		chain.IsValid(),
	)

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

		fmt.Println("Transactions:")

		for _, tx := range block.Transactions {
			fmt.Printf(
				"  %s -> %s : %.2f coins\n",
				tx.Sender,
				tx.Receiver,
				tx.Amount,
			)

			fmt.Printf(
				"  ID: %s\n",
				tx.ID,
			)
		}
	}
}

func openChain() *blockchain.Blockchain {
	chain, err := blockchain.OpenBlockchain(
		databasePath,
	)

	if err != nil {
		log.Fatal(err)
	}

	return chain
}

func closeChain(
	chain *blockchain.Blockchain,
) {
	if err := chain.Close(); err != nil {
		log.Printf(
			"failed to close blockchain: %v",
			err,
		)
	}
}
