package peer

import (
	"marabu/internal/crypto"
	"marabu/internal/messages"

	"fmt"
)

func (p *Peer) ValidateObject(obj Object) (ErrorCode, error) {
	objID, err := crypto.HashObject(obj)
	if err != nil {
		return messages.E_INTERNAL_ERROR, fmt.Errorf("Failed to hash object for validation: %v", err)
	}
	switch o := obj.(type) {
	case *messages.Transaction:
		fee, errorCode, err := p.ValidateTransaction(o)
		if err != nil {
			return errorCode, fmt.Errorf("Validation failed for transaction %s: %v", objID, err)
		}
		p.log(messages.MSG_OBJECT, fmt.Sprintf("Transaction %s is valid with fee %d", objID, fee))
		return E_NONE, nil
	case *messages.CoinbaseTransaction:
		fee, errorCode, err := p.ValidateCoinbase(o)
		if err != nil {
			return errorCode, fmt.Errorf("Validation failed for coinbase transaction %s: %v", objID, err)
		}
		p.log(messages.MSG_OBJECT, fmt.Sprintf("Coinbase transaction %s is valid with fee %d", objID, fee))
		return E_NONE, nil
	case *messages.Block:
		errorCode, err := p.ValidateBlock(o)
		if err != nil {
			return errorCode, fmt.Errorf("Validation failed for block %s: %v", objID, err)
		}
		p.log(messages.MSG_OBJECT, fmt.Sprintf("Block %s is valid", objID))
		return E_NONE, nil
	default:
		return messages.E_INTERNAL_ERROR, fmt.Errorf("Unknown object type: %T", obj)
	}
}

func (p *Peer) ValidateTransaction(tx *Transaction) (int, ErrorCode, error) {
	if tx.Type != messages.TRANSACTION {
		return 0, messages.E_INTERNAL_ERROR, fmt.Errorf("Invalid object type for transaction: %s", tx.Type)
	}

	sumInputs, sumOutputs := 0, 0

	for i, input := range tx.Inputs {
		outpoint := input.Outpoint
		// Find the referenced transaction
		exists, err := p.objectManager.Exists(outpoint.Txid)
		if !exists || err != nil {
			err = fmt.Errorf("Referenced transaction %s for input %d does not exist: %v", outpoint.Txid, i, err)
			return 0, messages.E_UNKNOWN_OBJECT, err
		}
		obj, err := p.objectManager.Get(outpoint.Txid)
		if err != nil {
			err = fmt.Errorf("Failed to fetch referenced transaction %s for input %d: %v", outpoint.Txid, i, err)
			return 0, messages.E_UNKNOWN_OBJECT, err
		}
		// Check that the referenced object is indeed a transaction or a coinbase transaction

		var outputs []messages.TxOutput
		// refTx, ok := obj.(*messages.Transaction)
		switch txObj := obj.(type) {
		case *Transaction:
			outputs = txObj.Outputs
		case *CoinbaseTransaction:
			outputs = txObj.Outputs
		case *Block:
			return 0, messages.E_INTERNAL_ERROR, fmt.Errorf("Referenced object %s for input %d is a block, expected transaction", outpoint.Txid, i)
		default:
			return 0, messages.E_INTERNAL_ERROR, fmt.Errorf("Referenced object %s for input %d is of unknown type, expected transaction", outpoint.Txid, i)
		}

		if outpoint.Index < 0 || outpoint.Index >= len(outputs) {
			return 0, messages.E_INVALID_TX_OUTPOINT, fmt.Errorf("Invalid output index %d in input %d referencing transaction %s", outpoint.Index, i, outpoint.Txid)
		}
		output := outputs[outpoint.Index]

		sumInputs += *output.Value

		if input.Sig == nil {
			return 0, E_INVALID_TX_SIGNATURE, fmt.Errorf("Missing signature for input %d referencing transaction %s", i, outpoint.Txid)
		}
		sig := string(*input.Sig)
		pubkey := string(output.Pubkey)
		msg := formatTransactionMessage(tx)

		if !crypto.Verify(pubkey, msg, sig) {
			err = fmt.Errorf("Invalid signature for input %d: pubkey %s, sig %s", i, pubkey, sig)
			return 0, E_INVALID_TX_SIGNATURE, err
		}
	}

	for _, output := range tx.Outputs {
		sumOutputs += *output.Value
	}

	if sumOutputs > sumInputs {
		return 0, E_INVALID_TX_CONSERVATION, fmt.Errorf("Output value %d exceeds input value %d", sumOutputs, sumInputs)
	}

	fee := sumInputs - sumOutputs
	return fee, E_NONE, nil
}

func (p *Peer) ValidateCoinbase(cb *CoinbaseTransaction) (int, ErrorCode, error) {
	return 0, E_NONE, nil
}

func (p *Peer) ValidateBlock(blk *Block) (ErrorCode, error) {
	return E_NONE, nil
}

// Creates a canonical message with nil signatures for signing/verification
func formatTransactionMessage(tx *Transaction) []byte {
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
