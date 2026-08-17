package blockchain

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// NodeServer exposes blockchain and peer operations over HTTP.
type NodeServer struct {
	Chain  *Blockchain
	Port   int
	Peers  *PeerSet
	Client *http.Client
}

// NewNodeServer creates a blockchain HTTP node.
func NewNodeServer(
	chain *Blockchain,
	port int,
) *NodeServer {
	return &NodeServer{
		Chain: chain,
		Port:  port,
		Peers: NewPeerSet(),
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Start starts the HTTP server.
func (server *NodeServer) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"/health",
		server.handleHealth,
	)

	mux.HandleFunc(
		"/chain",
		server.handleChain,
	)

	mux.HandleFunc(
		"/peers",
		server.handlePeers,
	)

	mux.HandleFunc(
		"/sync",
		server.handleSync,
	)

	mux.HandleFunc(
		"/fund",
		server.handleFund,
	)

	address := fmt.Sprintf(
		":%d",
		server.Port,
	)

	fmt.Printf(
		"ChainForge node listening on %s\n",
		address,
	)

	return http.ListenAndServe(
		address,
		mux,
	)
}

// handleHealth reports node status.
func (server *NodeServer) handleHealth(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodGet {
		http.Error(
			writer,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	writeJSON(
		writer,
		http.StatusOK,
		map[string]any{
			"status": "ok",
			"blocks": len(server.Chain.Blocks),
			"peers":  len(server.Peers.List()),
		},
	)
}

// handleChain returns the complete local blockchain.
func (server *NodeServer) handleChain(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodGet {
		http.Error(
			writer,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	writeJSON(
		writer,
		http.StatusOK,
		map[string]any{
			"length": len(server.Chain.Blocks),
			"valid":  server.Chain.IsValid(),
			"blocks": server.Chain.Blocks,
		},
	)
}

// handlePeers manages the local peer registry.
//
// GET  /peers
// POST /peers
func (server *NodeServer) handlePeers(
	writer http.ResponseWriter,
	request *http.Request,
) {
	switch request.Method {
	case http.MethodGet:
		writeJSON(
			writer,
			http.StatusOK,
			map[string]any{
				"peers": server.Peers.List(),
			},
		)

	case http.MethodPost:
		var payload struct {
			Peer string `json:"peer"`
		}

		if err := json.NewDecoder(
			request.Body,
		).Decode(&payload); err != nil {
			http.Error(
				writer,
				"invalid JSON",
				http.StatusBadRequest,
			)
			return
		}

		if err := server.Peers.Add(
			payload.Peer,
		); err != nil {
			http.Error(
				writer,
				err.Error(),
				http.StatusBadRequest,
			)
			return
		}

		writeJSON(
			writer,
			http.StatusCreated,
			map[string]any{
				"peer": payload.Peer,
			},
		)

	default:
		http.Error(
			writer,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

// handleFund creates a network-issued coinbase
// transaction and mines it directly on this node.
//
// POST /fund
func (server *NodeServer) handleFund(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodPost {
		http.Error(
			writer,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	var payload struct {
		Address string  `json:"address"`
		Amount  float64 `json:"amount"`
	}

	if err := json.NewDecoder(
		request.Body,
	).Decode(&payload); err != nil {
		http.Error(
			writer,
			"invalid JSON",
			http.StatusBadRequest,
		)
		return
	}

	if payload.Address == "" {
		http.Error(
			writer,
			"address is required",
			http.StatusBadRequest,
		)
		return
	}

	if payload.Amount <= 0 {
		http.Error(
			writer,
			"amount must be greater than zero",
			http.StatusBadRequest,
		)
		return
	}

	tx, err := NewCoinbaseTransaction(
		payload.Address,
		payload.Amount,
	)

	if err != nil {
		http.Error(
			writer,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	if err := server.Chain.AddBlock(
		[]Transaction{
			tx,
		},
	); err != nil {
		http.Error(
			writer,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(
		writer,
		http.StatusCreated,
		map[string]any{
			"status":  "funded",
			"address": payload.Address,
			"amount":  payload.Amount,
			"balance": server.Chain.BalanceOf(
				payload.Address,
			),
			"blocks": len(
				server.Chain.Blocks,
			),
		},
	)
}

// handleSync checks known peers and adopts
// the longest valid chain.
//
// POST /sync
func (server *NodeServer) handleSync(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodPost {
		http.Error(
			writer,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	updated, peer, err := server.Sync()

	if err != nil {
		http.Error(
			writer,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(
		writer,
		http.StatusOK,
		map[string]any{
			"updated": updated,
			"peer":    peer,
			"blocks": len(
				server.Chain.Blocks,
			),
		},
	)
}

// Sync checks all known peers and replaces the
// local chain with the longest valid chain found.
func (server *NodeServer) Sync() (
	bool,
	string,
	error,
) {
	bestBlocks := server.Chain.Blocks
	bestPeer := ""

	for _, peer := range server.Peers.List() {
		blocks, err := server.fetchPeerChain(
			peer,
		)

		if err != nil {
			continue
		}

		// Ignore chains that are not longer.
		if len(blocks) <= len(bestBlocks) {
			continue
		}

		candidate := &Blockchain{
			Blocks: blocks,
		}

		if !candidate.IsValid() {
			continue
		}

		bestBlocks = blocks
		bestPeer = peer
	}

	if bestPeer == "" {
		return false, "", nil
	}

	if err := server.replaceChain(
		bestBlocks,
	); err != nil {
		return false, "", err
	}

	return true, bestPeer, nil
}

// fetchPeerChain downloads a peer's blockchain.
func (server *NodeServer) fetchPeerChain(
	peer string,
) ([]*Block, error) {
	response, err := server.Client.Get(
		peer + "/chain",
	)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(
			"peer returned non-200 response",
		)
	}

	var payload struct {
		Length int      `json:"length"`
		Valid  bool     `json:"valid"`
		Blocks []*Block `json:"blocks"`
	}

	if err := json.NewDecoder(
		response.Body,
	).Decode(&payload); err != nil {
		return nil, err
	}

	if !payload.Valid {
		return nil, errors.New(
			"peer chain is invalid",
		)
	}

	if payload.Length != len(payload.Blocks) {
		return nil, errors.New(
			"peer chain length mismatch",
		)
	}

	return payload.Blocks, nil
}

// replaceChain stores and activates a synchronized chain.
func (server *NodeServer) replaceChain(
	blocks []*Block,
) error {
	if len(blocks) == 0 {
		return errors.New(
			"cannot replace chain with empty chain",
		)
	}

	candidate := &Blockchain{
		Blocks: blocks,
	}

	if !candidate.IsValid() {
		return errors.New(
			"replacement chain is invalid",
		)
	}

	// Persist synchronized blocks when this node
	// is using persistent storage.
	if server.Chain.Storage != nil {
		for _, block := range blocks {
			if err := server.Chain.Storage.SaveBlock(
				block,
			); err != nil {
				return err
			}
		}
	}

	server.Chain.Blocks = blocks

	return nil
}

// writeJSON writes a JSON HTTP response.
func writeJSON(
	writer http.ResponseWriter,
	status int,
	value any,
) {
	writer.Header().Set(
		"Content-Type",
		"application/json",
	)

	writer.WriteHeader(status)

	if err := json.NewEncoder(
		writer,
	).Encode(value); err != nil {
		fmt.Printf(
			"failed to encode response: %v\n",
			err,
		)
	}
}
