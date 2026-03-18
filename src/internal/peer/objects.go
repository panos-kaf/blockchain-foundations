package peer

import (
	"marabu/internal/crypto"
	"marabu/internal/messages"

	"fmt"
)

func (p *Peer) ValidateObject(obj messages.Object, objectID HashID) (messages.ErrorCode, error) {
	switch o := obj.(type) {
	case *messages.Transaction:
		fee, errorCode, err := p.ValidateTransaction(o, objectID)
		if err != nil {
			return errorCode, fmt.Errorf("Transaction validation failed: %v", err)
		}
		p.log(messages.OBJECT, fmt.Sprintf("Transaction %s is valid with fee %d", objectID, fee))
		return "", nil
	case *messages.CoinbaseTransaction:
		fee, errorCode, err := p.ValidateCoinbase(o, objectID)
		if err != nil {
			return errorCode, fmt.Errorf("Coinbase transaction validation failed: %v", err)
		}
		p.log(messages.OBJECT, fmt.Sprintf("Coinbase transaction %s is valid with fee %d", objectID, fee))
		return "", nil
	case *messages.Block:
		return p.ValidateBlock(o, objectID)
	default:
		return messages.INTERNAL_ERROR, fmt.Errorf("Unknown object type: %T", obj)
	}
}

func (p *Peer) ValidateTransaction(tx *messages.Transaction, objectID HashID) (int, messages.ErrorCode, error) {
	if tx.Type != messages.TRANSACTION {
		return 0, messages.INTERNAL_ERROR, fmt.Errorf("Invalid object type for transaction: %s", tx.Type)
	}

	sumInputs, sumOutputs := 0, 0

	for i, input := range tx.Inputs {
		outpoint := input.Outpoint
		// Find the referenced transaction
		obj, err := p.objectManager.FindObject(outpoint.Txid)
		if err != nil {
			err = fmt.Errorf("Failed to find referenced transaction %s for input %d: %v", outpoint.Txid, i, err)
			return 0, messages.UNKNOWN_OBJECT, err
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
			return 0, messages.INTERNAL_ERROR, fmt.Errorf("Referenced object %s for input %d is a block, expected transaction", outpoint.Txid, i)
		default:
			return 0, messages.INTERNAL_ERROR, fmt.Errorf("Referenced object %s for input %d is of unknown type, expected transaction", outpoint.Txid, i)
		}

		if outpoint.Index < 0 || outpoint.Index >= len(outputs) {
			return 0, messages.INVALID_TX_OUTPOINT, fmt.Errorf("Invalid output index %d in input %d referencing transaction %s", outpoint.Index, i, outpoint.Txid)
		}
		output := outputs[outpoint.Index]

		sumInputs += output.Value

		if input.Sig == nil {
			return 0, messages.INVALID_TX_SIGNATURE, fmt.Errorf("Missing signature for input %d referencing transaction %s", i, outpoint.Txid)
		}
		sig := string(*input.Sig)
		pubkey := string(output.Pubkey)
		msg := formatTransactionMessage(tx)

		if !crypto.Verify(pubkey, msg, sig) {
			err = fmt.Errorf("Invalid signature for input %d: pubkey %s, sig %s", i, pubkey, sig)
			return 0, messages.INVALID_TX_SIGNATURE, err
		}
	}

	for _, output := range tx.Outputs {
		sumOutputs += output.Value
	}

	if sumOutputs > sumInputs {
		return 0, messages.INVALID_TX_CONSERVATION, fmt.Errorf("Output value %d exceeds input value %d", sumOutputs, sumInputs)
	}

	fee := sumInputs - sumOutputs
	return fee, "", nil
}

func (p *Peer) ValidateCoinbase(cb *messages.CoinbaseTransaction, objectID HashID) (int, messages.ErrorCode, error) {
	return 0, "", nil
}

func (p *Peer) ValidateBlock(blk *messages.Block, objectID HashID) (messages.ErrorCode, error) {
	return "", nil
}

// Creates a canonical message with nil signatures for signing/verification
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
