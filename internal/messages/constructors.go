package messages

import (
	"encoding/json"
	"marabu/internal/crypto"
)

// -- Constructor functions for messages --

func MakeHelloMessage(version string, agent *string) (string, error) {
	return Canonicalize(HelloSchema{
		Type:    HELLO,
		Version: version,
		Agent:   agent,
	})
}

func MakeErrorMessage(name ErrorCode, description string) (string, error) {
	return Canonicalize(ErrorSchema{
		Type:        ERROR,
		Name:        name,
		Description: description,
	})
}

func MakeGetPeersMessage() (string, error) {
	return Canonicalize(GetPeersSchema{
		Type: GETPEERS,
	})
}

func MakePeersMessage(peers []string) (string, error) {
	return Canonicalize(PeersSchema{
		Type:  PEERS,
		Peers: peers,
	})
}

func MakeGetObjectMessage(objectID HashID) (string, error) {
	return Canonicalize(ObjectSchema{
		Type:     GETOBJECT,
		ObjectID: objectID,
	})
}

func MakeIHaveObjectMessage(objectID HashID) (string, error) {
	return Canonicalize(ObjectSchema{
		Type:     IHAVEOBJECT,
		ObjectID: objectID,
	})
}

func MakeTXObjectMessage(tx Transaction) (string, error) {
	raw, err := Canonicalize(tx)
	if err != nil {
		return "", err
	}
	objectID, err := crypto.HashString(raw)
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
	return Canonicalize(ObjectSchema{
		Type:      OBJECT,
		ObjectID:  objectID,
		RawObject: rawObject,
	})
}

func MakeGetMempoolMessage() (string, error) {
	return Canonicalize(GetMempoolSchema{
		Type: GETMEMPOOL,
	})
}

func MakeMempoolMessage(Txids []HashID) (string, error) {
	return Canonicalize(MempoolSchema{
		Type:  MEMPOOL,
		Txids: Txids,
	})
}

func MakeGetChainTipMessage() (string, error) {
	return Canonicalize(GetChainTipSchema{
		Type: GETCHAINTIP,
	})
}

func MakeChainTipMessage(BlockID HashID) (string, error) {
	return Canonicalize(ChainTipSchema{
		Type:  CHAINTIP,
		Block: BlockID,
	})
}
