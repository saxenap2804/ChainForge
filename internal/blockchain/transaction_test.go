package blockchain

import "testing"

func TestValidTransaction(
	t *testing.T,
) {
	tx, err := NewTransaction(
		"Alice",
		"Bob",
		10,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if tx.ID == "" {
		t.Fatal(
			"expected transaction ID",
		)
	}
}

func TestTransactionRejectsEmptySender(
	t *testing.T,
) {
	_, err := NewTransaction(
		"",
		"Bob",
		10,
	)

	if err == nil {
		t.Fatal(
			"expected empty sender to be rejected",
		)
	}
}

func TestTransactionRejectsEmptyReceiver(
	t *testing.T,
) {
	_, err := NewTransaction(
		"Alice",
		"",
		10,
	)

	if err == nil {
		t.Fatal(
			"expected empty receiver to be rejected",
		)
	}
}

func TestTransactionRejectsNegativeAmount(
	t *testing.T,
) {
	_, err := NewTransaction(
		"Alice",
		"Bob",
		-5,
	)

	if err == nil {
		t.Fatal(
			"expected negative amount to be rejected",
		)
	}
}

func TestTransactionRejectsSelfTransfer(
	t *testing.T,
) {
	_, err := NewTransaction(
		"Alice",
		"Alice",
		10,
	)

	if err == nil {
		t.Fatal(
			"expected self-transfer to be rejected",
		)
	}
}

func TestTransactionIDChangesWithContents(
	t *testing.T,
) {
	tx1, err := NewTransaction(
		"Alice",
		"Bob",
		10,
	)

	if err != nil {
		t.Fatal(err)
	}

	tx2, err := NewTransaction(
		"Alice",
		"Charlie",
		10,
	)

	if err != nil {
		t.Fatal(err)
	}

	if tx1.ID == tx2.ID {
		t.Fatal(
			"expected different transactions to have different IDs",
		)
	}
}

func TestTransactionTrimsAddresses(
	t *testing.T,
) {
	tx, err := NewTransaction(
		"  Alice  ",
		"  Bob  ",
		10,
	)

	if err != nil {
		t.Fatal(err)
	}

	if tx.Sender != "Alice" {
		t.Fatalf(
			"expected sender Alice, got %q",
			tx.Sender,
		)
	}

	if tx.Receiver != "Bob" {
		t.Fatalf(
			"expected receiver Bob, got %q",
			tx.Receiver,
		)
	}
}

func TestTransactionIDIsDeterministic(
	t *testing.T,
) {
	tx1, err := NewTransaction(
		"Alice",
		"Bob",
		10,
	)

	if err != nil {
		t.Fatal(err)
	}

	tx2, err := NewTransaction(
		"Alice",
		"Bob",
		10,
	)

	if err != nil {
		t.Fatal(err)
	}

	if tx1.ID != tx2.ID {
		t.Fatal(
			"identical transactions should have identical IDs",
		)
	}
}
