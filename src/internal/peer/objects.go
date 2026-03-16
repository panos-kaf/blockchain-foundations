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
		refTx, ok := obj.(*messages.Transaction)
		if !ok {
			return 0, fmt.Errorf("Referenced object %s for input %d is not a transaction", outpoint.Txid, i)
		}
		if outpoint.Index < 0 || outpoint.Index >= len(refTx.Outputs) {
			return 0, fmt.Errorf("Invalid output index %d in input %d referencing transaction %s", outpoint.Index, i, outpoint.Txid)
		}
		output := refTx.Outputs[outpoint.Index]
		sig := input.Sig

		sumInputs += output.Value

		pubkey, err := crypto.StringToPubkey(string(output.Pubkey))
		if err != nil {
			return 0, fmt.Errorf("Invalid pubkey hex string in input %d: %v", i, err)
		}

		if !crypto.Verify([]byte(objectID), string(sig), pubkey) {
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

func (p *Peer) ValidateCoinbase(cb *messages.CoinbaseTransaction, objectID HashID) (int, error) {
	return 0, nil
}

func (p *Peer) ValidateBlock(blk *messages.Block, objectID HashID) error {
	return nil
}
