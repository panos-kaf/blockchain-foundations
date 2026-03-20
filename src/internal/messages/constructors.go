package messages

import (
	"encoding/json"
)

// -- Constructor functions for messages --

func MakeHelloMessage() (string, error) {

	version := T_Version("0.10.0")
	agent := T_BuString("marabobos")

	return CanonicalizeMessage(&HelloMessage{
		Type:      MSG_HELLO,
		T_Version: version,
		Agent:     &agent,
	})
}

func MakeErrorMessage(name ErrorCode, description T_BuString) (string, error) {
	return CanonicalizeMessage(&ErrorMessage{
		Type:        MSG_ERROR,
		Name:        name,
		Description: description,
	})
}

func MakeGetPeersMessage() (string, error) {
	return CanonicalizeMessage(&GetPeersMessage{
		Type: MSG_GETPEERS,
	})
}

func MakePeersMessage(peers T_Peers) (string, error) {
	return CanonicalizeMessage(&PeersMessage{
		Type:    MSG_PEERS,
		T_Peers: peers,
	})
}

func MakeGetObjectMessage(objectID T_HashID) (string, error) {
	return CanonicalizeMessage(&GetObjectMessage{
		Type:     MSG_GETOBJECT,
		ObjectID: objectID,
	})
}

func MakeIHaveObjectMessage(objectID T_HashID) (string, error) {
	return CanonicalizeMessage(&IHaveObjectMessage{
		Type:     MSG_IHAVEOBJECT,
		ObjectID: objectID,
	})
}

func MakeObjectMessage(obj Object) (string, error) {
	raw, err := Canonicalize(obj)
	if err != nil {
		return "", err
	}
	return CanonicalizeMessage(&ObjectMessage{
		Type:      MSG_OBJECT,
		RawObject: json.RawMessage(raw),
	})
}

func MakeGetMempoolMessage() (string, error) {
	return CanonicalizeMessage(&GetMempoolMessage{
		Type: MSG_GETMEMPOOL,
	})
}

func MakeMempoolMessage(Txids T_HashIDs) (string, error) {
	return CanonicalizeMessage(&MempoolMessage{
		Type:  MSG_MEMPOOL,
		Txids: Txids,
	})
}

func MakeGetChainTipMessage() (string, error) {
	return CanonicalizeMessage(&GetChainTipMessage{
		Type: MSG_GETCHAINTIP,
	})
}

func MakeChainTipMessage(BlockID T_HashID) (string, error) {
	return CanonicalizeMessage(&ChainTipMessage{
		Type:    MSG_CHAINTIP,
		T_Block: BlockID,
	})
}
