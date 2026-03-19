package messages

import (
	"encoding/json"
	"fmt"
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
		Type:    MSG_HELLO,
		Version: version,
		Agent:   &agent,
	})
}

func MakeErrorMessage(name ErrorCode, description string) (string, error) {
	return CanonicalizeMessage(ErrorSchema{
		Type:        MSG_ERROR,
		Name:        name,
		Description: description,
	})
}

func MakeGetPeersMessage() (string, error) {
	return CanonicalizeMessage(GetPeersSchema{
		Type: MSG_GETPEERS,
	})
}

func MakePeersMessage(peers []string) (string, error) {
	return CanonicalizeMessage(PeersSchema{
		Type:  MSG_PEERS,
		Peers: peers,
	})
}

func MakeGetObjectMessage(objectID HashID) (string, error) {
	return CanonicalizeMessage(GetObjectSchema{
		Type:     MSG_GETOBJECT,
		ObjectID: objectID,
	})
}

func MakeIHaveObjectMessage(objectID HashID) (string, error) {
	return CanonicalizeMessage(IHaveObjectSchema{
		Type:     MSG_IHAVEOBJECT,
		ObjectID: objectID,
	})
}

func MakeObjectMessage(obj Object) (string, error) {
	raw, err := Canonicalize(obj)
	if err != nil {
		return "", err
	}
	return CanonicalizeMessage(ObjectSchema{
		Type:      MSG_OBJECT,
		RawObject: json.RawMessage(raw),
	})
}

func MakeGetMempoolMessage() (string, error) {
	return CanonicalizeMessage(GetMempoolSchema{
		Type: MSG_GETMEMPOOL,
	})
}

func MakeMempoolMessage(Txids []HashID) (string, error) {
	return CanonicalizeMessage(MempoolSchema{
		Type:  MSG_MEMPOOL,
		Txids: Txids,
	})
}

func MakeGetChainTipMessage() (string, error) {
	return CanonicalizeMessage(GetChainTipSchema{
		Type: MSG_GETCHAINTIP,
	})
}

func MakeChainTipMessage(BlockID HashID) (string, error) {
	return CanonicalizeMessage(ChainTipSchema{
		Type:  MSG_CHAINTIP,
		Block: BlockID,
	})
}
