package peer

import (
	"marabu/internal/crypto"
	"marabu/internal/messages"

	"fmt"
)

func (p *Peer) ValidateObject(obj Object) (ErrorCode, error) {
	objID, err := crypto.HashObject(obj)
	if err != nil {
		return E_INTERNAL_ERROR, fmt.Errorf("Failed to hash object for validation: %v", err)
	}
	switch o := obj.(type) {
	case *Transaction:
		fee, errorCode, err := p.ValidateTransaction(o)
		if err != nil {
			return errorCode, fmt.Errorf("Validation failed for transaction %s: %v", objID, err)
		}
		p.log(MSG_OBJECT, fmt.Sprintf("Transaction %s is valid with fee %d", objID, fee))
		return E_NONE, nil
	case *CoinbaseTransaction:
		fee, errorCode, err := p.ValidateCoinbase(o)
		if err != nil {
			return errorCode, fmt.Errorf("Validation failed for coinbase transaction %s: %v", objID, err)
		}
		p.log(MSG_OBJECT, fmt.Sprintf("Coinbase transaction %s is valid with fee %d", objID, fee))
		return E_NONE, nil
	case *Block:
		errorCode, err := p.ValidateBlock(o)
		if err != nil {
			return errorCode, fmt.Errorf("Validation failed for block %s: %v", objID, err)
		}
		p.log(MSG_OBJECT, fmt.Sprintf("Block %s is valid", objID))
		return E_NONE, nil
	default:
		return E_INTERNAL_ERROR, fmt.Errorf("Unknown object type: %T", obj)
	}
}

func (p *Peer) ValidateTransaction(tx *Transaction) (int, ErrorCode, error) {
	if tx.Type != messages.OBJ_TRANSACTION {
		return 0, messages.E_INTERNAL_ERROR, fmt.Errorf("Invalid object type for transaction: %s", tx.Type)
	}

	sumInputs, sumOutputs := 0, 0

	// 1. Create a quick struct to hold the data we need for crypto later
	type sigData struct {
		pubkey string
		sig    string
	}
	var verifyQueue []sigData

	// input/output transaction validity checks
	for i, input := range tx.Inputs {
		outpoint := input.Outpoint

		exists, err := p.objectManager.Exists(outpoint.Txid)
		if !exists || err != nil {
			return 0, messages.E_UNKNOWN_OBJECT, fmt.Errorf("Referenced transaction %s for input %d does not exist", outpoint.Txid, i)
		}

		obj, err := p.objectManager.Get(outpoint.Txid)
		if err != nil {
			return 0, messages.E_UNKNOWN_OBJECT, fmt.Errorf("Failed to fetch referenced transaction")
		}

		var outputs []messages.TxOutput
		switch txObj := obj.(type) {
		case *Transaction:
			outputs = txObj.Outputs
		case *CoinbaseTransaction:
			outputs = txObj.Outputs
		default:
			return 0, messages.E_INTERNAL_ERROR, fmt.Errorf("Referenced object is of unknown type")
		}

		if outpoint.Index < 0 || outpoint.Index >= len(outputs) {
			return 0, messages.E_INVALID_TX_OUTPOINT, fmt.Errorf("Invalid output index")
		}

		output := outputs[outpoint.Index]
		sumInputs += *output.Value

		if input.Sig == nil {
			return 0, E_INVALID_TX_SIGNATURE, fmt.Errorf("Missing signature")
		}

		// cache pubkey and sig for later verification
		verifyQueue = append(verifyQueue, sigData{
			pubkey: string(output.Pubkey),
			sig:    string(*input.Sig),
		})
	}

	// conservation check
	for _, output := range tx.Outputs {
		sumOutputs += *output.Value
	}

	if sumOutputs > sumInputs {
		return 0, E_INVALID_TX_CONSERVATION, fmt.Errorf("Output value %d exceeds input value %d", sumOutputs, sumInputs)
	}

	// sig verification
	msg := messages.TxMessageForSignature(tx)

	for i, data := range verifyQueue {
		if !crypto.Verify(data.pubkey, msg, data.sig) {
			return 0, E_INVALID_TX_SIGNATURE, fmt.Errorf("Invalid signature for input %d", i)
		}
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
