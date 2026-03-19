package messages

import (
	"encoding/json"
)

type (
	Message interface {
		MessageType() MessageType
		Validate() (error, ErrorCode)
	}

	MessageType string
	ErrorCode   string
	HashID      string // 32byte (64-character) hex string
	Signature   string // 64byte (128-character) hex string
)

const (
	MSG_NONE        MessageType = ""
	MSG_HELLO       MessageType = "hello"
	MSG_ERROR       MessageType = "error"
	MSG_GETPEERS    MessageType = "getpeers"
	MSG_PEERS       MessageType = "peers"
	MSG_GETOBJECT   MessageType = "getobject"
	MSG_IHAVEOBJECT MessageType = "ihaveobject"
	MSG_OBJECT      MessageType = "object"
	MSG_GETMEMPOOL  MessageType = "getmempool"
	MSG_MEMPOOL     MessageType = "mempool"
	MSG_GETCHAINTIP MessageType = "getchaintip"
	MSG_CHAINTIP    MessageType = "chaintip"

	E_NONE                    ErrorCode = ""
	E_INTERNAL_ERROR          ErrorCode = "E_INTERNAL_ERROR"
	E_INVALID_FORMAT          ErrorCode = "E_INVALID_FORMAT"
	E_UNKNOWN_OBJECT          ErrorCode = "E_UNKNOWN_OBJECT"
	E_UNFINDABLE_OBJECT       ErrorCode = "E_UNFINDABLE_OBJECT"
	E_INVALID_HANDSHAKE       ErrorCode = "E_INVALID_HANDSHAKE"
	E_INVALID_TX_OUTPOINT     ErrorCode = "E_INVALID_TX_OUTPOINT"
	E_INVALID_TX_SIGNATURE    ErrorCode = "E_INVALID_TX_SIGNATURE"
	E_INVALID_TX_CONSERVATION ErrorCode = "E_INVALID_TX_CONSERVATION"
	E_INVALID_BLOCK_COINBASE  ErrorCode = "E_INVALID_BLOCK_COINBASE"
	E_INVALID_BLOCK_TIMESTAMP ErrorCode = "E_INVALID_BLOCK_TIMESTAMP"
	E_INVALID_BLOCK_POW       ErrorCode = "E_INVALID_BLOCK_POW"
	E_INVALID_GENESIS         ErrorCode = "E_INVALID_GENESIS"
)

type (
	HelloSchema struct {
		Type    MessageType `json:"type"`
		Version string      `json:"version"`
		Agent   *string     `json:"agent,omitempty"`
	}

	ErrorSchema struct {
		Type        MessageType `json:"type"`
		Name        ErrorCode   `json:"name"`
		Description string      `json:"description"`
	}

	GetPeersSchema struct {
		Type MessageType `json:"type"`
	}

	PeersSchema struct {
		Type  MessageType `json:"type"`
		Peers []string    `json:"peers"`
	}

	GetObjectSchema struct {
		Type     MessageType `json:"type"`
		ObjectID HashID      `json:"objectid"`
	}

	IHaveObjectSchema struct {
		Type     MessageType `json:"type"`
		ObjectID HashID      `json:"objectid"`
	}

	ObjectSchema struct {
		Type MessageType `json:"type"`

		// The raw, unparsed JSON of the object.
		RawObject json.RawMessage `json:"object"`

		// This field is not part of the JSON schema but is used internally to hold the deserialized object after validation
		Object Object `json:"-"`
	}

	GetMempoolSchema struct {
		Type MessageType `json:"type"`
	}

	MempoolSchema struct {
		Type  MessageType `json:"type"`
		Txids []HashID    `json:"txids"`
	}

	GetChainTipSchema struct {
		Type MessageType `json:"type"`
	}

	ChainTipSchema struct {
		Type  MessageType `json:"type"`
		Block HashID      `json:"block"`
	}
)

// -- message type getters --

func (h *HelloSchema) MessageType() MessageType {
	return MSG_HELLO
}
func (e *ErrorSchema) MessageType() MessageType {
	return MSG_ERROR
}
func (g *GetPeersSchema) MessageType() MessageType {
	return MSG_GETPEERS
}
func (p *PeersSchema) MessageType() MessageType {
	return MSG_PEERS
}
func (g *GetObjectSchema) MessageType() MessageType {
	return MSG_GETOBJECT
}
func (i *IHaveObjectSchema) MessageType() MessageType {
	return MSG_IHAVEOBJECT
}
func (o *ObjectSchema) MessageType() MessageType {
	return MSG_OBJECT
}
func (g *GetMempoolSchema) MessageType() MessageType {
	return MSG_GETMEMPOOL
}
func (m *MempoolSchema) MessageType() MessageType {
	return MSG_MEMPOOL
}
func (g *GetChainTipSchema) MessageType() MessageType {
	return MSG_GETCHAINTIP
}
func (c *ChainTipSchema) MessageType() MessageType {
	return MSG_CHAINTIP
}
