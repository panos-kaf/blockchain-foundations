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
	Sig Signature `json:"sig"`
}

type TxOutput struct {
	Pubkey HashID `json:"pubkey"`
	Value  int    `json:"value"`
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
	Height  int        `json:"height"`
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

// tx or coinbase or block
type Object interface {
	ObjectType() ObjectType
	Validate() error
}

func (t *Transaction) ObjectType() ObjectType {
	return TRANSACTION
}

func (c *CoinbaseTransaction) ObjectType() ObjectType {
	return TRANSACTION
}

func (b *Block) ObjectType() ObjectType {
	return BLOCK
}

func (t *Transaction) Validate() error {

	if err := ValidateObjectType(t.Type); err != nil {
		return err
	}
	arrLength := len(t.Inputs)
	if arrLength == 0 {
		return fmt.Errorf("transaction must have at least one input")
	}
	if arrLength > 1000 {
		return fmt.Errorf("transaction exceeds maximum number of inputs (1000), got %d", arrLength)
	}
	for i, input := range t.Inputs {
		if err := ValidateNonNegativeInt(input.Outpoint.Index, fmt.Sprintf("inputs[%d].outpoint.index", i)); err != nil {
			return err
		}
	}

	arrLength = len(t.Outputs)
	if arrLength > 1000 {
		return fmt.Errorf("transaction exceeds maximum number of outputs (1000), got %d", arrLength)
	}
	for i, output := range t.Outputs {
		if err := ValidateNonNegativeInt(output.Value, fmt.Sprintf("outputs[%d].value", i)); err != nil {
			return err
		}
	}
	return nil
}

func (c *CoinbaseTransaction) Validate() error {
	return nil
}

func (b *Block) Validate() error {
	return nil
}
