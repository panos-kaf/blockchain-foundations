package messages

import (
	"encoding/json"
	"fmt"
	"marabu/internal/crypto"
)

// wraps the canonicalization process for message constructors
func CanonicalizeMessage(msg interface{}) (string, error) {

	canon, err := Canonicalize(msg)
	if err != nil {
		return "", fmt.Errorf("Error parsing message: %w", err)
	}
	return canon + "\n", nil
}

// -- Constructor functions for messages --

func MakeHelloMessage() (string, error) {

	version := "0.10.0"
	agent := "marabobos"

	return CanonicalizeMessage(HelloSchema{
		Type:    HELLO,
		Version: version,
		Agent:   &agent,
	})
}

func MakeErrorMessage(name ErrorCode, description string) (string, error) {
	return CanonicalizeMessage(ErrorSchema{
		Type:        ERROR,
		Name:        name,
		Description: description,
	})
}

func MakeGetPeersMessage() (string, error) {
	return CanonicalizeMessage(GetPeersSchema{
		Type: GETPEERS,
	})
}

func MakePeersMessage(peers []string) (string, error) {
	return CanonicalizeMessage(PeersSchema{
		Type:  PEERS,
		Peers: peers,
	})
}

func MakeGetObjectMessage(objectID HashID) (string, error) {
	return CanonicalizeMessage(GetObjectSchema{
		Type:     GETOBJECT,
		ObjectID: objectID,
	})
}

func MakeIHaveObjectMessage(objectID HashID) (string, error) {
	return CanonicalizeMessage(IHaveObjectSchema{
		Type:     IHAVEOBJECT,
		ObjectID: objectID,
	})
}

func MakeTXObjectMessage(tx Transaction) (string, error) {
	raw, err := Canonicalize(tx)
	if err != nil {
		return "", err
	}
	objectID, err := HashObject(tx)
	if err != nil {
		return "", err
	}
	return MakeObjectMessage(HashID(objectID), json.RawMessage(raw))
}

func MakeCBTXObjectMessage(cbtx CoinbaseTransaction) (string, error) {
	raw, err := Canonicalize(cbtx)
	if err != nil {
		return "", err
	}
	objectID, err := crypto.HashString(raw)
	if err != nil {
		return "", err
	}

	return MakeObjectMessage(HashID(objectID), json.RawMessage(raw))
}

func MakeBlockObjectMessage(block Block) (string, error) {
	raw, err := Canonicalize(block)
	if err != nil {
		return "", err
	}
	objectID, err := crypto.HashString(raw)
	if err != nil {
		return "", err
	}

	return MakeObjectMessage(HashID(objectID), json.RawMessage(raw))
}

func MakeObjectMessage(objectID HashID, rawObject json.RawMessage) (string, error) {
	return CanonicalizeMessage(ObjectSchema{
		Type:      OBJECT,
		ObjectID:  objectID,
		RawObject: rawObject,
	})
}

func MakeGetMempoolMessage() (string, error) {
	return CanonicalizeMessage(GetMempoolSchema{
		Type: GETMEMPOOL,
	})
}

func MakeMempoolMessage(Txids []HashID) (string, error) {
	return CanonicalizeMessage(MempoolSchema{
		Type:  MEMPOOL,
		Txids: Txids,
	})
}

func MakeGetChainTipMessage() (string, error) {
	return CanonicalizeMessage(GetChainTipSchema{
		Type: GETCHAINTIP,
	})
}

func MakeChainTipMessage(BlockID HashID) (string, error) {
	return CanonicalizeMessage(ChainTipSchema{
		Type:  CHAINTIP,
		Block: BlockID,
	})
}
