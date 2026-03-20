package messages

import (
	"fmt"
)

// -- Object sub-type definitions --

type (
	T_Outpoint struct {
		Txid  T_HashID `json:"txid"`
		Index T_BuInt  `json:"index"`
	}

	T_TxInput struct {
		T_Outpoint T_Outpoint `json:"outpoint"`

		// 64byte (128-character) hexadecimal string, handle as simple string for now...
		Sig *T_Signature `json:"sig"`
	}

	T_TxOutput struct {
		Pubkey T_HashID `json:"pubkey"`
		Value  *T_BuInt `json:"value"`
	}

	ObjectType string
)

const (
	OBJ_BLOCK       ObjectType = "block"
	OBJ_TRANSACTION ObjectType = "transaction"
)

type T_Transaction struct {
	Type    ObjectType   `json:"type"`
	Inputs  []T_TxInput  `json:"inputs"`
	Outputs []T_TxOutput `json:"outputs"`
}

type T_CoinbaseTransaction struct {
	Type    ObjectType   `json:"type"`
	Height  *T_BuInt     `json:"height"`
	Outputs []T_TxOutput `json:"outputs"`
}

type T_Block struct {
	Type       ObjectType   `json:"type"`
	T          T_HashID     `json:"T"`
	Created    T_BuInt      `json:"created"`
	Miner      *T_BuString  `json:"miner,omitempty"`
	Nonce      T_HashID     `json:"nonce"`
	Note       *T_BuString  `json:"note,omitempty"`
	Previd     *T_HashID    `json:"previd"` //nullable
	Studentids *T_BuStrings `json:"studentids,omitempty"`
	Txids      T_HashIDs    `json:"txids"`
}

// An object is either a Tx, a coinbase Tx, or a block
type Object interface {
	ObjectType() ObjectType
	Validate() (error, ErrorCode)
}

func (t T_Transaction) ObjectType() ObjectType {
	return OBJ_TRANSACTION
}

func (c T_CoinbaseTransaction) ObjectType() ObjectType {
	return OBJ_TRANSACTION
}

func (b T_Block) ObjectType() ObjectType {
	return OBJ_BLOCK
}

func (t T_Transaction) Validate() (error, ErrorCode) {

	arrLength := len(t.Inputs)
	if arrLength == 0 {
		return fmt.Errorf("transaction must have at least one input"), E_INVALID_FORMAT
	}
	if arrLength > 1000 {
		return fmt.Errorf("transaction exceeds maximum number of inputs (1000), got %d", arrLength), E_INVALID_FORMAT
	}

	for i, output := range t.Outputs {
		if output.Value == nil {
			return fmt.Errorf("missing value for output %d", i), E_INVALID_FORMAT
		}
	}

	arrLength = len(t.Outputs)
	if arrLength > 1000 {
		return fmt.Errorf("transaction exceeds maximum number of outputs (1000), got %d", arrLength), E_INVALID_FORMAT
	}

	return nil, E_NONE
}

func (c T_CoinbaseTransaction) Validate() (error, ErrorCode) {

	if c.Height == nil {
		return fmt.Errorf("missing height for coinbase transaction"), E_INVALID_FORMAT
	}

	for i, output := range c.Outputs {
		if output.Value == nil {
			return fmt.Errorf("missing value for output %d", i), E_INVALID_FORMAT
		}
	}
	return nil, E_NONE
}

func (b T_Block) Validate() (error, ErrorCode) {
	return nil, E_NONE
}

// -- object constructors --

func makeTxInput(txid T_HashID, index T_BuInt, sig T_Signature) T_TxInput {
	return T_TxInput{
		T_Outpoint: T_Outpoint{
			Txid:  txid,
			Index: index,
		},
		Sig: &sig,
	}
}

func makeTxOutput(pubkey T_HashID, value T_BuInt) T_TxOutput {
	return T_TxOutput{
		Pubkey: pubkey,
		Value:  &value,
	}
}

func makeTransaction(inputs []T_TxInput, outputs []T_TxOutput) T_Transaction {
	return T_Transaction{
		Type:    OBJ_TRANSACTION,
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func makeCoinbaseTransaction(height T_BuInt, outputs []T_TxOutput) T_CoinbaseTransaction {
	return T_CoinbaseTransaction{
		Type:    OBJ_TRANSACTION,
		Height:  &height,
		Outputs: outputs,
	}
}

func makeBlock(T T_HashID, created T_BuInt, miner *T_BuString, nonce T_HashID, note *T_BuString, previd *T_HashID, studentids *T_BuStrings, txids []T_HashID) T_Block {
	return T_Block{
		Type:       OBJ_BLOCK,
		T:          T,
		Created:    created,
		Miner:      miner,
		Nonce:      nonce,
		Note:       note,
		Previd:     previd,
		Studentids: studentids,
		Txids:      txids,
	}
}
