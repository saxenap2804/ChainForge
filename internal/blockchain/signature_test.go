package blockchain

import "testing"

func TestTransactionSignatureValid(
	t *testing.T,
) {
	wallet, err := NewWallet()

	if err != nil {
		t.Fatal(err)
	}

	receiver, err := NewWallet()

	if err != nil {
		t.Fatal(err)
	}

	tx, err := NewTransaction(
		wallet.Address,
		receiver.Address,
		10,
	)

	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Sign(wallet); err != nil {
		t.Fatal(err)
	}

	if !tx.Verify() {
		t.Fatal(
			"expected signed transaction to verify",
		)
	}
}

func TestTransactionSignatureFailsAfterTampering(
	t *testing.T,
) {
	wallet, err := NewWallet()

	if err != nil {
		t.Fatal(err)
	}

	receiver, err := NewWallet()

	if err != nil {
		t.Fatal(err)
	}

	tx, err := NewTransaction(
		wallet.Address,
		receiver.Address,
		10,
	)

	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Sign(wallet); err != nil {
		t.Fatal(err)
	}

	tx.Amount = 1000

	if tx.Verify() {
		t.Fatal(
			"expected tampered transaction signature to fail",
		)
	}
}

func TestWrongWalletCannotSignTransaction(
	t *testing.T,
) {
	alice, err := NewWallet()

	if err != nil {
		t.Fatal(err)
	}

	bob, err := NewWallet()

	if err != nil {
		t.Fatal(err)
	}

	receiver, err := NewWallet()

	if err != nil {
		t.Fatal(err)
	}

	tx, err := NewTransaction(
		alice.Address,
		receiver.Address,
		5,
	)

	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Sign(bob); err == nil {
		t.Fatal(
			"expected wrong wallet signing attempt to fail",
		)
	}
}
