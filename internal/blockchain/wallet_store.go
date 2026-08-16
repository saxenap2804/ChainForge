package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"errors"
	"math/big"
	"os"
)

type storedWallet struct {
	D         []byte `json:"d"`
	PublicKey []byte `json:"public_key"`
	Address   string `json:"address"`
}

type WalletStore struct {
	path string
}

// NewWalletStore creates a wallet store backed by a JSON file.
func NewWalletStore(path string) *WalletStore {
	return &WalletStore{
		path: path,
	}
}

// Load loads all wallets from disk.
func (ws *WalletStore) Load() (map[string]*Wallet, error) {
	wallets := make(map[string]*Wallet)

	data, err := os.ReadFile(ws.path)

	if errors.Is(err, os.ErrNotExist) {
		return wallets, nil
	}

	if err != nil {
		return nil, err
	}

	var stored map[string]storedWallet

	if err := json.Unmarshal(data, &stored); err != nil {
		return nil, err
	}

	for address, sw := range stored {
		curve := elliptic.P256()

		privateKey := new(ecdsa.PrivateKey)
		privateKey.PublicKey.Curve = curve
		privateKey.D = newBigInt(sw.D)

		x, y := elliptic.Unmarshal(
			curve,
			sw.PublicKey,
		)

		if x == nil || y == nil {
			return nil, errors.New("invalid stored public key")
		}

		privateKey.PublicKey.X = x
		privateKey.PublicKey.Y = y

		wallets[address] = &Wallet{
			PrivateKey: privateKey,
			PublicKey:  sw.PublicKey,
			Address:    sw.Address,
		}
	}

	return wallets, nil
}

// Save writes all wallets to disk.
func (ws *WalletStore) Save(
	wallets map[string]*Wallet,
) error {
	stored := make(
		map[string]storedWallet,
	)

	for address, wallet := range wallets {
		stored[address] = storedWallet{
			D:         wallet.PrivateKey.D.Bytes(),
			PublicKey: wallet.PublicKey,
			Address:   wallet.Address,
		}
	}

	data, err := json.MarshalIndent(
		stored,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		ws.path,
		data,
		0600,
	)
}

func newBigInt(data []byte) *big.Int {
	return new(big.Int).SetBytes(data)
}
