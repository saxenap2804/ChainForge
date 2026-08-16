package blockchain

import "testing"

func TestNewWallet(t *testing.T) {
	wallet, err := NewWallet()

	if err != nil {
		t.Fatalf(
			"unexpected wallet error: %v",
			err,
		)
	}

	if wallet.PrivateKey == nil {
		t.Fatal(
			"expected private key",
		)
	}

	if len(wallet.PublicKey) == 0 {
		t.Fatal(
			"expected public key",
		)
	}

	if wallet.Address == "" {
		t.Fatal(
			"expected wallet address",
		)
	}
}

func TestWalletsGenerateDifferentAddresses(
	t *testing.T,
) {
	walletOne, err := NewWallet()

	if err != nil {
		t.Fatal(err)
	}

	walletTwo, err := NewWallet()

	if err != nil {
		t.Fatal(err)
	}

	if walletOne.Address == walletTwo.Address {
		t.Fatal(
			"expected different wallets to have different addresses",
		)
	}
}
