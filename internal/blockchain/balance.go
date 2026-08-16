package blockchain

import "errors"

// BalanceOf calculates the current balance
// for a wallet address by scanning the chain.
func (bc *Blockchain) BalanceOf(address string) float64 {
	var balance float64

	for _, block := range bc.Blocks {
		for _, tx := range block.Transactions {
			if tx.Receiver == address {
				balance += tx.Amount
			}

			if tx.Sender == address {
				balance -= tx.Amount
			}
		}
	}

	return balance
}

// ValidateTransaction verifies transaction structure,
// signature, and available sender balance.
func (bc *Blockchain) ValidateTransaction(
	tx Transaction,
) error {
	if err := tx.Validate(); err != nil {
		return err
	}

	// Coinbase transactions create new coins.
	if tx.IsCoinbase() {
		return nil
	}

	if !tx.Verify() {
		return errors.New(
			"invalid transaction signature",
		)
	}

	currentBalance := bc.BalanceOf(
		tx.Sender,
	)

	if currentBalance < tx.Amount {
		return errors.New(
			"insufficient balance",
		)
	}

	return nil
}
