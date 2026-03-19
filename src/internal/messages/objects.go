package messages

import (
	"fmt"
)

// -- Object sub-type definitions --

type Outpoint struct {
	Txid  HashID `json:"txid"`
	Index int    `json:"index"`
}

type TxInput struct {
	Outpoint Outpoint `json:"outpoint"`

	// 64byte (128-character) hexadecimal string, handle as simple string for now...
	Sig *Signature `json:"sig"`
}

type TxOutput struct {
	Pubkey HashID `json:"pubkey"`
	Value  *int   `json:"value"`
}

type ObjectType string

const (
	BLOCK       ObjectType = "block"
	TRANSACTION ObjectType = "transaction"
)

type Transaction struct {
	Type    ObjectType `json:"type"`
	Inputs  []TxInput  `json:"inputs"`
	Outputs []TxOutput `json:"outputs"`
}

type CoinbaseTransaction struct {
	Type    ObjectType `json:"type"`
	Height  *int       `json:"height"`
	Outputs []TxOutput `json:"outputs"`
}

type Block struct {
	Type       ObjectType `json:"type"`
	T          HashID     `json:"T"`
	Created    int        `json:"created"`
	Miner      *string    `json:"miner,omitempty"`
	Nonce      HashID     `json:"nonce"`
	Note       *string    `json:"note,omitempty"`
	Previd     *HashID    `json:"previd"` //nullable
	Studentids *[]string  `json:"studentids,omitempty"`
	Txids      []string   `json:"txids"`
}

// An object is either a Tx, a coinbase Tx, or a block
type Object interface {
	ObjectType() ObjectType
	Validate() (error, ErrorCode)
}

func (t Transaction) ObjectType() ObjectType {
	return TRANSACTION
}

func (c CoinbaseTransaction) ObjectType() ObjectType {
	return TRANSACTION
}

func (b Block) ObjectType() ObjectType {
	return BLOCK
}

func (t Transaction) Validate() (error, ErrorCode) {

	if err, code := ValidateObjectType(t.Type); err != nil {
		return err, code
	}
	arrLength := len(t.Inputs)
	if arrLength == 0 {
		return fmt.Errorf("transaction must have at least one input"), INVALID_FORMAT
	}
	if arrLength > 1000 {
		return fmt.Errorf("transaction exceeds maximum number of inputs (1000), got %d", arrLength), INVALID_FORMAT
	}
	for i, input := range t.Inputs {
		if err, code := ValidateNonNegativeInt(input.Outpoint.Index, fmt.Sprintf("inputs[%d].outpoint.index", i)); err != nil {
			return err, code
		}
	}
	for i, output := range t.Outputs {
		if output.Value == nil {
			return fmt.Errorf("missing value for output %d", i), INVALID_FORMAT
		}
	}

	arrLength = len(t.Outputs)
	if arrLength > 1000 {
		return fmt.Errorf("transaction exceeds maximum number of outputs (1000), got %d", arrLength), INVALID_FORMAT
	}
	for i, output := range t.Outputs {
		if err, code := ValidateNonNegativeInt(*output.Value, fmt.Sprintf("outputs[%d].value", i)); err != nil {
			return err, code
		}
	}
	return nil, ""
}

func (c CoinbaseTransaction) Validate() (error, ErrorCode) {

	if err, code := ValidateObjectType(c.Type); err != nil {
		return err, code
	}

	if c.Height == nil {
		return fmt.Errorf("missing height for coinbase transaction"), INVALID_FORMAT
	}

	for i, output := range c.Outputs {
		if output.Value == nil {
			return fmt.Errorf("missing value for output %d", i), INVALID_FORMAT
		}
	}
	return nil, ""
}

func (b Block) Validate() (error, ErrorCode) {
	return nil, ""
}

// -- object constructors --

func makeTxInput(txid HashID, index int, sig Signature) TxInput {
	return TxInput{
		Outpoint: Outpoint{
			Txid:  txid,
			Index: index,
		},
		Sig: &sig,
	}
}

func makeTxOutput(pubkey HashID, value int) TxOutput {
	return TxOutput{
		Pubkey: pubkey,
		Value:  &value,
	}
}

func makeTransaction(inputs []TxInput, outputs []TxOutput) Transaction {
	return Transaction{
		Type:    TRANSACTION,
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func makeCoinbaseTransaction(height int, outputs []TxOutput) CoinbaseTransaction {
	return CoinbaseTransaction{
		Type:    TRANSACTION,
		Height:  &height,
		Outputs: outputs,
	}
}

func makeBlock(T HashID, created int, miner *string, nonce HashID, note *string, previd *HashID, studentids *[]string, txids []string) Block {
	return Block{
		Type:       BLOCK,
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
