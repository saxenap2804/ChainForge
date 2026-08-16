package blockchain

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// NodeServer exposes blockchain data over HTTP.
type NodeServer struct {
	Chain *Blockchain
	Port  int
}

// NewNodeServer creates a new HTTP blockchain node.
func NewNodeServer(
	chain *Blockchain,
	port int,
) *NodeServer {
	return &NodeServer{
		Chain: chain,
		Port:  port,
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

// handleHealth reports whether the node is online.
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

	writer.Header().Set(
		"Content-Type",
		"application/json",
	)

	response := map[string]any{
		"status": "ok",
		"blocks": len(server.Chain.Blocks),
	}

	_ = json.NewEncoder(
		writer,
	).Encode(response)
}

// handleChain returns the complete blockchain.
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

	writer.Header().Set(
		"Content-Type",
		"application/json",
	)

	response := map[string]any{
		"length": len(server.Chain.Blocks),
		"valid":  server.Chain.IsValid(),
		"blocks": server.Chain.Blocks,
	}

	if err := json.NewEncoder(
		writer,
	).Encode(response); err != nil {
		http.Error(
			writer,
			"failed to encode blockchain",
			http.StatusInternalServerError,
		)
	}
}
