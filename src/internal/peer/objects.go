package peer

import (
	"marabu/internal/crypto"
	"marabu/internal/messages"

	"fmt"
)

func (p *Peer) ValidateObject(obj messages.Object, objectID HashID) error {
	switch o := obj.(type) {
	case *messages.Transaction:
		fee, err := p.ValidateTransaction(o, objectID)
		if err != nil {
			return fmt.Errorf("Transaction validation failed: %v", err)
		}
		p.log(messages.OBJECT, fmt.Sprintf("Transaction %s is valid with fee %d", objectID, fee))
		return nil
	case *messages.CoinbaseTransaction:
		fee, err := p.ValidateCoinbase(o, objectID)
		if err != nil {
			return fmt.Errorf("Coinbase transaction validation failed: %v", err)
		}
		p.log(messages.OBJECT, fmt.Sprintf("Coinbase transaction %s is valid with fee %d", objectID, fee))
		return nil
	case *messages.Block:
		return p.ValidateBlock(o, objectID)
	default:
		return fmt.Errorf("Unknown object type: %T", obj)
	}
}

func (p *Peer) ValidateTransaction(tx *messages.Transaction, objectID HashID) (int, error) {
	if tx.Type != messages.TRANSACTION {
		return 0, fmt.Errorf("Invalid object type for transaction: %s", tx.Type)
	}

	sumInputs, sumOutputs := 0, 0

	for i, input := range tx.Inputs {
		outpoint := input.Outpoint
		// Find the referenced transaction
		obj, err := p.objectManager.FindObject(outpoint.Txid)
		if err != nil {
			return 0, fmt.Errorf("Failed to find referenced transaction %s for input %d: %v", outpoint.Txid, i, err)
		}
		// Check that the referenced object is indeed a transaction or a coinbase transaction

		var outputs []messages.TxOutput
		// refTx, ok := obj.(*messages.Transaction)
		switch txObj := obj.(type) {
		case *messages.Transaction:
			outputs = txObj.Outputs
		case *messages.CoinbaseTransaction:
			outputs = txObj.Outputs
		case *messages.Block:
			return 0, fmt.Errorf("Referenced object %s for input %d is a block, expected transaction", outpoint.Txid, i)
		default:
			return 0, fmt.Errorf("Referenced object %s for input %d is of unknown type, expected transaction", outpoint.Txid, i)
		}

		if outpoint.Index < 0 || outpoint.Index >= len(outputs) {
			return 0, fmt.Errorf("Invalid output index %d in input %d referencing transaction %s", outpoint.Index, i, outpoint.Txid)
		}
		output := outputs[outpoint.Index]

		sumInputs += output.Value

		sig := string(*input.Sig)
		pubkey := string(output.Pubkey)
		msg := formatTransactionMessage(tx)

		if !crypto.Verify(pubkey, msg, sig) {
			return 0, fmt.Errorf("Invalid signature for input %d: does not match referenced output's pubkey", i)
		}
	}

	for _, output := range tx.Outputs {
		sumOutputs += output.Value
	}

	if sumOutputs > sumInputs {
		return 0, fmt.Errorf("Output value %d exceeds input value %d", sumOutputs, sumInputs)
	}

	fee := sumInputs - sumOutputs
	return fee, nil

}

func formatTransactionMessage(tx *messages.Transaction) []byte {
	// Create a copy of the transaction with empty signatures for signing/verification
	txCopy := *tx
	txCopy.Inputs = make([]messages.TxInput, len(tx.Inputs))
	copy(txCopy.Inputs, tx.Inputs)
	for i := range txCopy.Inputs {
		txCopy.Inputs[i].Sig = nil
	}
	// Canonicalize the transaction copy to get the message bytes
	msgBytes, _ := (messages.Canonicalize(txCopy))
	return []byte(msgBytes)
}

func (p *Peer) ValidateCoinbase(cb *messages.CoinbaseTransaction, objectID HashID) (int, error) {
	return 0, nil
}

func (p *Peer) ValidateBlock(blk *messages.Block, objectID HashID) error {
	return nil
}
