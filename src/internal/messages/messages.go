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
	HELLO       MessageType = "hello"
	ERROR       MessageType = "error"
	GETPEERS    MessageType = "getpeers"
	PEERS       MessageType = "peers"
	GETOBJECT   MessageType = "getobject"
	IHAVEOBJECT MessageType = "ihaveobject"
	OBJECT      MessageType = "object"
	GETMEMPOOL  MessageType = "getmempool"
	MEMPOOL     MessageType = "mempool"
	GETCHAINTIP MessageType = "getchaintip"
	CHAINTIP    MessageType = "chaintip"

	INTERNAL_ERROR          ErrorCode = "INTERNAL_ERROR"
	INVALID_FORMAT          ErrorCode = "INVALID_FORMAT"
	UNKNOWN_OBJECT          ErrorCode = "UNKNOWN_OBJECT"
	UNFINDABLE_OBJECT       ErrorCode = "UNFINDABLE_OBJECT"
	INVALID_HANDSHAKE       ErrorCode = "INVALID_HANDSHAKE"
	INVALID_TX_OUTPOINT     ErrorCode = "INVALID_TX_OUTPOINT"
	INVALID_TX_SIGNATURE    ErrorCode = "INVALID_TX_SIGNATURE"
	INVALID_TX_CONSERVATION ErrorCode = "INVALID_TX_CONSERVATION"
	INVALID_BLOCK_COINBASE  ErrorCode = "INVALID_BLOCK_COINBASE"
	INVALID_BLOCK_TIMESTAMP ErrorCode = "INVALID_BLOCK_TIMESTAMP"
	INVALID_BLOCK_POW       ErrorCode = "INVALID_BLOCK_POW"
	INVALID_GENESIS         ErrorCode = "INVALID_GENESIS"
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
	return HELLO
}
func (e *ErrorSchema) MessageType() MessageType {
	return ERROR
}
func (g *GetPeersSchema) MessageType() MessageType {
	return GETPEERS
}
func (p *PeersSchema) MessageType() MessageType {
	return PEERS
}
func (g *GetObjectSchema) MessageType() MessageType {
	return GETOBJECT
}
func (i *IHaveObjectSchema) MessageType() MessageType {
	return IHAVEOBJECT
}
func (o *ObjectSchema) MessageType() MessageType {
	return OBJECT
}
func (g *GetMempoolSchema) MessageType() MessageType {
	return GETMEMPOOL
}
func (m *MempoolSchema) MessageType() MessageType {
	return MEMPOOL
}
func (g *GetChainTipSchema) MessageType() MessageType {
	return GETCHAINTIP
}
func (c *ChainTipSchema) MessageType() MessageType {
	return CHAINTIP
}
