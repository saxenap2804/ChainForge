package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

	case "sendnode":
		sendNodeTransaction(os.Args[2:])

	case "printchain":
		printBlockchain()

	case "startnode":
		startNode(os.Args[2:])

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
	fmt.Println("  sendnode --node URL --from ADDRESS --to ADDRESS --amount AMOUNT")
	fmt.Println("  printchain")
	fmt.Println("  startnode --port PORT --db DATABASE")
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
// of an address from the default local chain.
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
// on the default local blockchain.
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
// mines, and persists a transaction
// directly to the default local blockchain.
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

// sendNodeTransaction creates and signs a transaction
// locally, then submits it to a running ChainForge node.
func sendNodeTransaction(args []string) {
	command := flag.NewFlagSet(
		"sendnode",
		flag.ExitOnError,
	)

	node := command.String(
		"node",
		"",
		"node URL, for example http://localhost:8080",
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

	if *node == "" {
		log.Fatal(
			"--node is required",
		)
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

	nodeURL := strings.TrimRight(
		*node,
		"/",
	)

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
			"sender wallet not found in wallets.json",
		)
	}

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

	payload, err := json.Marshal(
		tx,
	)

	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	response, err := client.Post(
		nodeURL+"/transactions",
		"application/json",
		bytes.NewReader(payload),
	)

	if err != nil {
		log.Fatalf(
			"failed to contact node: %v",
			err,
		)
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(
		response.Body,
	)

	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode != http.StatusCreated {
		log.Fatalf(
			"node rejected transaction (%s): %s",
			response.Status,
			strings.TrimSpace(
				string(responseBody),
			),
		)
	}

	var result struct {
		Status string  `json:"status"`
		ID     string  `json:"id"`
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
		Blocks int     `json:"blocks"`
	}

	if len(responseBody) > 0 {
		if err := json.Unmarshal(
			responseBody,
			&result,
		); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(
		"Transaction submitted successfully.",
	)

	fmt.Printf(
		"Node: %s\n",
		nodeURL,
	)

	fmt.Printf(
		"%s -> %s : %.2f coins\n",
		*from,
		*to,
		*amount,
	)

	fmt.Printf(
		"Transaction ID: %s\n",
		tx.ID,
	)

	if result.Blocks > 0 {
		fmt.Printf(
			"Node chain length: %d\n",
			result.Blocks,
		)
	}
}

// printBlockchain displays every block
// stored in the default local blockchain.
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

		fmt.Println(
			"Transactions:",
		)

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

// startNode starts an HTTP node using
// the specified blockchain database.
// startNode starts an HTTP node using
// the specified blockchain database.
//
// Peers can be configured automatically with
// CHAINFORGE_PEERS as a comma-separated list.
//
// Example:
//
// CHAINFORGE_PEERS=http://node-b:8081,http://node-c:8082
func startNode(args []string) {
	command := flag.NewFlagSet(
		"startnode",
		flag.ExitOnError,
	)

	port := command.Int(
		"port",
		8080,
		"HTTP node port",
	)

	dbPath := command.String(
		"db",
		databasePath,
		"blockchain database path",
	)

	if err := command.Parse(args); err != nil {
		log.Fatal(err)
	}

	if *port <= 0 || *port > 65535 {
		log.Fatal(
			"port must be between 1 and 65535",
		)
	}

	chain, err := blockchain.OpenBlockchain(
		*dbPath,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer closeChain(chain)

	server := blockchain.NewNodeServer(
		chain,
		*port,
	)

	// ------------------------------------------------
	// Automatically register peers from environment.
	// ------------------------------------------------

	peerEnv := strings.TrimSpace(
		os.Getenv("CHAINFORGE_PEERS"),
	)

	if peerEnv != "" {
		for _, peer := range strings.Split(
			peerEnv,
			",",
		) {
			peer = strings.TrimSpace(peer)

			if peer == "" {
				continue
			}

			if err := server.Peers.Add(peer); err != nil {
				log.Printf(
					"Could not register peer %s: %v",
					peer,
					err,
				)

				continue
			}

			log.Printf(
				"Registered peer: %s",
				peer,
			)
		}
	}

	// ------------------------------------------------
	// Initial synchronization.
	//
	// Docker services may not become available at
	// exactly the same moment, so retry briefly.
	// ------------------------------------------------

	if len(server.Peers.List()) > 0 {
		go func() {
			for attempt := 1; attempt <= 5; attempt++ {
				time.Sleep(2 * time.Second)

				updated, peer, err := server.Sync()

				if err != nil {
					log.Printf(
						"Initial sync attempt %d failed: %v",
						attempt,
						err,
					)

					continue
				}

				if updated {
					log.Printf(
						"Blockchain synchronized from %s",
						peer,
					)
				}

				log.Printf(
					"Sync attempt %d complete. Blocks: %d",
					attempt,
					len(server.Chain.Blocks),
				)
			}
		}()
	}

	if err := server.Start(); err != nil {
		log.Fatal(err)
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
	if chain == nil {
		return
	}

	if err := chain.Close(); err != nil {
		log.Printf(
			"failed to close blockchain: %v",
			err,
		)
	}
}
