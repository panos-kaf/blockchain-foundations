package peer

import (
	"marabu/internal/messages"
)

const (
	sent = true
	recv = false
)

var (
	MSG_NONE        = messages.MSG_NONE
	MSG_HELLO       = messages.MSG_HELLO
	MSG_ERROR       = messages.MSG_ERROR
	MSG_GETPEERS    = messages.MSG_GETPEERS
	MSG_PEERS       = messages.MSG_PEERS
	MSG_GETOBJECT   = messages.MSG_GETOBJECT
	MSG_IHAVEOBJECT = messages.MSG_IHAVEOBJECT
	MSG_OBJECT      = messages.MSG_OBJECT
	MSG_GETMEMPOOL  = messages.MSG_GETMEMPOOL
	MSG_MEMPOOL     = messages.MSG_MEMPOOL
	MSG_GETCHAINTIP = messages.MSG_GETCHAINTIP
	MSG_CHAINTIP    = messages.MSG_CHAINTIP

	E_NONE                    = messages.E_NONE
	E_INTERNAL_ERROR          = messages.E_INTERNAL_ERROR
	E_INVALID_FORMAT          = messages.E_INVALID_FORMAT
	E_UNKNOWN_OBJECT          = messages.E_UNKNOWN_OBJECT
	E_UNFINDABLE_OBJECT       = messages.E_UNFINDABLE_OBJECT
	E_INVALID_HANDSHAKE       = messages.E_INVALID_HANDSHAKE
	E_INVALID_TX_OUTPOINT     = messages.E_INVALID_TX_OUTPOINT
	E_INVALID_TX_SIGNATURE    = messages.E_INVALID_TX_SIGNATURE
	E_INVALID_TX_CONSERVATION = messages.E_INVALID_TX_CONSERVATION
	E_INVALID_BLOCK_COINBASE  = messages.E_INVALID_BLOCK_COINBASE
	E_INVALID_BLOCK_TIMESTAMP = messages.E_INVALID_BLOCK_TIMESTAMP
	E_INVALID_BLOCK_POW       = messages.E_INVALID_BLOCK_POW
	E_INVALID_GENESIS         = messages.E_INVALID_GENESIS
)

type (
	Message     = messages.Message
	MessageType = messages.MessageType
	ErrorCode   = messages.ErrorCode

	HelloSchema       = messages.HelloSchema
	ErrorSchema       = messages.ErrorSchema
	GetPeersSchema    = messages.GetPeersSchema
	PeersSchema       = messages.PeersSchema
	GetObjectSchema   = messages.GetObjectSchema
	IHaveObjectSchema = messages.IHaveObjectSchema
	ObjectSchema      = messages.ObjectSchema
	GetMempoolSchema  = messages.GetMempoolSchema
	MempoolSchema     = messages.MempoolSchema
	GetChainTipSchema = messages.GetChainTipSchema
	ChainTipSchema    = messages.ChainTipSchema

	Object              = messages.Object
	Transaction         = messages.Transaction
	CoinbaseTransaction = messages.CoinbaseTransaction
	Block               = messages.Block

	HashID    = messages.HashID
	Signature = messages.Signature
)
