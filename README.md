# ChainForge

ChainForge is a distributed blockchain network built from scratch in Go.

It implements the core components of a blockchain system, including Proof of Work mining, cryptographically signed transactions, wallet generation, balance validation, persistent storage, REST APIs, peer-to-peer synchronization, and multi-node deployment with Docker.

The project was built to explore how blockchain systems work below the framework level—from block hashing and transaction validation to distributed node communication and chain synchronization.

---

## Features

- SHA-256 block hashing
- Deterministic genesis block
- Proof of Work mining
- Nonce-based block validation
- Blockchain integrity verification
- Tamper detection
- Structured transactions
- Deterministic transaction IDs
- ECDSA wallet generation
- Digital transaction signatures
- Signature verification
- Balance tracking
- Overspending protection
- Coinbase/funding transactions
- Persistent blockchain storage with BoltDB
- Persistent wallet storage
- Command-line interface
- REST API
- Peer discovery and registration
- Blockchain synchronization
- Block propagation between nodes
- Multi-node Docker deployment
- Docker volumes for persistent node data
- Health checks
- Automated Go test suite
- GitHub Actions CI

---

## Architecture

```text
                         ChainForge
                             |
            +----------------+----------------+
            |                                 |
        Blockchain                          Wallets
            |                                 |
     +------+------+                   +------+------+
     |             |                   |             |
   Blocks     Transactions           ECDSA       Addresses
     |             |
     |        +----+----+
     |        |         |
     |    Signatures  Balances
     |
 Proof of Work
     |
 SHA-256 Mining
     |
 Persistent Storage
     |
   BoltDB
     |
 REST / P2P Layer
     |
 +---+-------------------+
 |                       |
Node A                  Node B
:8080                   :8081
 |                       |
 +------ Block Sync -----+
```

Each node maintains its own persistent blockchain database and communicates with peers through HTTP.

---

## Distributed Network

The included Docker Compose configuration launches two independent ChainForge nodes.

```text
+------------------------+
|        Node A          |
|                        |
| HTTP: localhost:8080   |
| DB: /data/node-a.db    |
+-----------+------------+
            |
            | block propagation
            | chain synchronization
            |
+-----------+------------+
|        Node B          |
|                        |
| HTTP: localhost:8081   |
| DB: /data/node-b.db    |
+------------------------+
```

Each node:

- maintains its own blockchain
- stores data independently
- tracks peers
- validates incoming blocks
- exposes REST endpoints
- synchronizes blockchain state
- persists data across container restarts

---

## Block Structure

A ChainForge block contains:

```go
type Block struct {
    Timestamp     int64
    Transactions  []Transaction
    PrevBlockHash []byte
    Hash          []byte
    Nonce         int64
}
```

Blocks are cryptographically linked through `PrevBlockHash`.

Changing transaction data in an existing block invalidates its Proof of Work and breaks blockchain integrity.

---

## Proof of Work

ChainForge implements a SHA-256 based Proof of Work algorithm.

For each block, the miner searches for a nonce such that:

```text
SHA256(
    previous block hash +
    transaction hashes +
    timestamp +
    difficulty +
    nonce
)
```

produces a value below the configured target.

Conceptually:

```text
hash < target
```

The resulting nonce and hash are stored in the block.

During blockchain validation, the Proof of Work is recalculated to ensure that the block has not been modified.

---

## Deterministic Genesis Block

All ChainForge nodes begin from the same deterministic genesis block.

This is important in a distributed blockchain because independently generated genesis blocks would have different hashes and therefore represent different chains.

A shared genesis block allows independently started nodes to validate and synchronize blocks from the same network.

---

## Transactions

Transactions represent transfers between blockchain addresses.

```go
type Transaction struct {
    ID        string
    Sender    string
    Receiver  string
    Amount    float64
    PublicKey []byte
    Signature []byte
}
```

Each transaction receives a deterministic SHA-256 identifier.

ChainForge validates transaction rules including:

- sender must exist
- receiver must exist
- sender and receiver cannot be identical
- transfer values must be valid
- sender must have sufficient balance
- signed transactions must pass cryptographic verification

---

## Wallets and Digital Signatures

ChainForge generates wallets using ECDSA with the P-256 elliptic curve.

Each wallet contains:

```text
Private Key
     |
     v
Public Key
     |
     v
Blockchain Address
```

Before a transaction is submitted, the sender signs the transaction using their private key.

The blockchain verifies the signature using the corresponding public key.

This prevents another wallet from authorizing transactions on behalf of the sender.

---

## Balance Validation

Balances are derived from blockchain transaction history.

For an address:

```text
Balance =
    incoming transactions
    - outgoing transactions
```

Before accepting a transaction, ChainForge verifies that the sender has enough funds.

Example:

```text
Alice balance: 50

Alice -> Bob: 10

Alice balance: 40
Bob balance:   10
```

An attempted transfer greater than the sender's available balance is rejected.

---

## Persistent Storage

ChainForge uses BoltDB (`bbolt`) for blockchain persistence.

Each node stores blockchain data in its own database.

Examples:

```text
chainforge.db
node-a.db
node-b.db
```

Docker nodes use persistent volumes:

```text
node-a-data
node-b-data
```

This allows blockchain state to survive process or container restarts.

---

## CLI

ChainForge provides a command-line interface for interacting with the blockchain.

### Create a wallet

```bash
go run . createwallet
```

Example:

```text
Wallet created successfully.
Address: cd57509d75b3ad744b6ea9bbcafce0820ac7bf5d
```

### Check balance

```bash
go run . balance --address ADDRESS
```

### Fund an address

```bash
go run . fund --address ADDRESS --amount 50
```

### Send coins

```bash
go run . send \
  --from SENDER_ADDRESS \
  --to RECEIVER_ADDRESS \
  --amount 10
```

### Print blockchain

```bash
go run . printchain
```

### Start a blockchain node

```bash
go run . startnode --port 8080 --db node-a.db
```

A second node can be started with:

```bash
go run . startnode --port 8081 --db node-b.db
```

### Submit a transaction to a running node

```bash
go run . sendnode \
  --node http://localhost:8080 \
  --from SENDER_ADDRESS \
  --to RECEIVER_ADDRESS \
  --amount 10
```

---

## REST API

Running nodes expose an HTTP API.

### Health

```http
GET /health
```

Example response:

```json
{
  "blocks": 2,
  "peers": 1,
  "status": "ok"
}
```

### Blockchain

```http
GET /chain
```

Returns the blockchain maintained by the node.

### Peers

```http
GET /peers
```

Returns registered peers.

A peer can be registered through:

```http
POST /peers
```

Example body:

```json
{
  "peer": "http://node-b:8081"
}
```

### Synchronize

```http
POST /sync
```

Requests blockchain synchronization with known peers.

### Fund Address

```http
POST /fund
```

Example:

```json
{
  "address": "WALLET_ADDRESS",
  "amount": 10
}
```

### Transaction Submission

ChainForge also supports submitting signed transactions to a running node through its transaction API, which is used by the `sendnode` CLI command.

---

## Peer Synchronization

Nodes maintain a list of known peers.

In Docker, peers communicate using service names:

```text
node-a:8080
node-b:8081
```

while the host machine accesses them through:

```text
localhost:8080
localhost:8081
```

When new blockchain state is produced, peer nodes can validate and synchronize the chain.

A successful propagation test results in:

```text
Node A blocks: 2
Node B blocks: 2
```

without requiring both nodes to share the same database.

---

## Running with Docker

### Requirements

- Docker Desktop
- Docker Compose

### Build and start

```bash
docker compose up --build -d
```

Check container status:

```bash
docker compose ps
```

Expected result:

```text
chainforge-node-a   Up (healthy)
chainforge-node-b   Up (healthy)
```

### Test Node A

```text
http://localhost:8080/health
```

### Test Node B

```text
http://localhost:8081/health
```

### Stop the network

```bash
docker compose down
```

To remove the node data volumes as well:

```bash
docker compose down -v
```

---

## Example Distributed Flow

```text
Create Wallet
     |
     v
Fund Wallet
     |
     v
Node A receives transaction
     |
     v
Validate transaction
     |
     v
Check sender balance
     |
     v
Verify ECDSA signature
     |
     v
Mine block using Proof of Work
     |
     v
Persist block to BoltDB
     |
     v
Broadcast / synchronize
     |
     v
Node B validates blockchain state
     |
     v
Both nodes converge
```

---

## Testing

Run the complete test suite with:

```bash
go test ./... -v
```

The tests cover areas including:

- genesis block creation
- blockchain initialization
- block creation
- blockchain validation
- tamper detection
- Proof of Work
- transaction validation
- deterministic transaction IDs
- transaction signatures
- invalid signature detection
- wallet generation
- balance behavior
- persistent blockchain reopening
- BoltDB storage

---

## Continuous Integration

ChainForge uses GitHub Actions to automatically validate pushes and pull requests.

The CI pipeline performs:

```text
Checkout
   |
   v
Set up Go
   |
   v
Download Dependencies
   |
   v
Formatting Check
   |
   v
Run Tests
   |
   v
Build ChainForge
```

Workflow:

```text
.github/workflows/tests.yml
```

---

## Project Structure

```text
ChainForge/
|
|-- .github/
|   `-- workflows/
|       `-- tests.yml
|
|-- internal/
|   `-- blockchain/
|       |-- balance.go
|       |-- block.go
|       |-- blockchain.go
|       |-- peer.go
|       |-- proof_of_work.go
|       |-- server.go
|       |-- storage.go
|       |-- transaction.go
|       |-- wallet.go
|       |-- wallet_store.go
|       |
|       |-- blockchain_test.go
|       |-- persistence_test.go
|       |-- proof_of_work_test.go
|       |-- signature_test.go
|       |-- storage_test.go
|       |-- transaction_test.go
|       `-- wallet_test.go
|
|-- .dockerignore
|-- .gitignore
|-- compose.yml
|-- Dockerfile
|-- go.mod
|-- go.sum
|-- LICENSE
|-- main.go
`-- README.md
```

---

## Technology Stack

| Area | Technology |
|---|---|
| Language | Go |
| Hashing | SHA-256 |
| Digital Signatures | ECDSA P-256 |
| Persistence | BoltDB / bbolt |
| Networking | Go `net/http` |
| API | REST |
| Containers | Docker |
| Multi-node orchestration | Docker Compose |
| Testing | Go `testing` |
| CI | GitHub Actions |

---

## Security and Design Concepts Demonstrated

ChainForge demonstrates several foundational distributed-system and blockchain concepts:

- cryptographic hashing
- immutable hash-linked data structures
- Proof of Work
- deterministic network genesis state
- asymmetric cryptography
- digital signatures
- transaction authorization
- state derived from transaction history
- persistent storage
- peer communication
- distributed state synchronization
- containerized multi-node deployment

ChainForge is an educational implementation intended to demonstrate blockchain and distributed-systems engineering concepts. It is not intended for production cryptocurrency or financial use.

---

## Future Improvements

Potential extensions include:

- Merkle trees
- UTXO-based transaction accounting
- transaction mempool
- configurable mining difficulty
- mining rewards
- richer peer discovery
- fork resolution and chain-selection policies
- more advanced consensus mechanisms
- authenticated node communication
- observability and metrics

---

## License

This project is licensed under the terms provided in the repository's `LICENSE` file.