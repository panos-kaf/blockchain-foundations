package messages

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var messageTypeRegistry = map[string]reflect.Type{
	string(MSG_HELLO):       reflect.TypeOf(HelloSchema{}),
	string(MSG_ERROR):       reflect.TypeOf(ErrorSchema{}),
	string(MSG_GETPEERS):    reflect.TypeOf(GetPeersSchema{}),
	string(MSG_PEERS):       reflect.TypeOf(PeersSchema{}),
	string(MSG_GETOBJECT):   reflect.TypeOf(GetObjectSchema{}),
	string(MSG_IHAVEOBJECT): reflect.TypeOf(IHaveObjectSchema{}),
	string(MSG_OBJECT):      reflect.TypeOf(ObjectSchema{}),
	string(MSG_GETMEMPOOL):  reflect.TypeOf(GetMempoolSchema{}),
	string(MSG_MEMPOOL):     reflect.TypeOf(MempoolSchema{}),
	string(MSG_GETCHAINTIP): reflect.TypeOf(GetChainTipSchema{}),
	string(MSG_CHAINTIP):    reflect.TypeOf(ChainTipSchema{}),
}

func UnmarshalMessage(raw string) (Message, error, ErrorCode) {
	var probe map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &probe); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err), E_INVALID_FORMAT
	}
	typeVal, ok := probe["type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'type' field in message"), E_INVALID_FORMAT
	}

	typ, found := messageTypeRegistry[typeVal]
	if !found {
		return nil, fmt.Errorf("unknown message type: '%s'", typeVal), E_INVALID_FORMAT
	}

	msgPtr := reflect.New(typ)
	if err := json.Unmarshal([]byte(raw), msgPtr.Interface()); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s message: %w", typeVal, err), E_INVALID_FORMAT
	}

	// Return as Message interface
	if m, ok := msgPtr.Interface().(Message); ok {
		return m, nil, E_NONE
	}
	return nil, fmt.Errorf("type %s does not implement Message interface", typeVal), E_INVALID_FORMAT
}

// Custom UnmarshalJSON for MessageType to enforce valid message types
func (mt *MessageType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch MessageType(s) {
	case MSG_HELLO, MSG_ERROR, MSG_GETPEERS, MSG_PEERS, MSG_GETOBJECT, MSG_IHAVEOBJECT, MSG_OBJECT, MSG_GETMEMPOOL, MSG_MEMPOOL, MSG_GETCHAINTIP, MSG_CHAINTIP:
		*mt = MessageType(s)
		return nil
	default:
		return fmt.Errorf("invalid message type: '%s'", s)
	}
}

// Custom UnmarshalJSON for ObjectType to enforce valid object types
func (ot *ObjectType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch ObjectType(s) {
	case OBJ_BLOCK, OBJ_TRANSACTION:
		*ot = ObjectType(s)
		return nil
	default:
		return fmt.Errorf("invalid object type: '%s'", s)
	}
}

// Custom UnmarshalJSON for ErrorCode to enforce valid error codes
func (ec *ErrorCode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch ErrorCode(s) {
	case E_INTERNAL_ERROR, E_INVALID_FORMAT, E_UNKNOWN_OBJECT, E_UNFINDABLE_OBJECT, E_INVALID_HANDSHAKE, E_INVALID_TX_OUTPOINT, E_INVALID_TX_SIGNATURE, E_INVALID_TX_CONSERVATION, E_INVALID_BLOCK_COINBASE, E_INVALID_BLOCK_TIMESTAMP, E_INVALID_BLOCK_POW, E_INVALID_GENESIS:
		*ec = ErrorCode(s)
		return nil
	default:
		return fmt.Errorf("invalid error code: '%s'", s)
	}
}

// Custom UnmarshalJSON for Hash IDs to enforce length and hex format
func (h *HashID) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if len(s) != 64 {
		return fmt.Errorf("invalid hash ID: must be exactly 64 characters, got %d", len(s))
	}

	// Hexification - Hex strings must be in lower case.
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return fmt.Errorf("invalid hash ID: must be hexadecimal, got invalid character '%c'", c)
		}
	}

	*h = HashID(s)
	return nil
}

// Custom UnmarshalJSON for Signatures to enforce length and hex format
func (sig *Signature) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if len(s) != 128 {
		return fmt.Errorf("invalid signature: must be exactly 128 characters, got %d", len(s))
	}

	// Hexification - Hex strings must be in lower case.
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return fmt.Errorf("invalid signature: must be hexadecimal, got invalid character '%c'", c)
		}
	}

	*sig = Signature(s)
	return nil
}

// Custom UnmarshalJSON for ObjectSchema to handle dynamic inner object types (block, transaction, coinbase transaction)
func (o *ObjectSchema) UnmarshalJSON(data []byte) error {

	type Alias ObjectSchema

	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var typeProbe struct {
		Type ObjectType `json:"type"`
	}
	if err := json.Unmarshal(o.RawObject, &typeProbe); err != nil {
		return fmt.Errorf("failed to probe inner object type: %w", err)
	}

	// populate the Object field based on the type of the inner object
	switch typeProbe.Type {
	case OBJ_BLOCK:
		var b Block
		if err := json.Unmarshal(o.RawObject, &b); err != nil {
			return fmt.Errorf("failed to unmarshal block object: %w", err)
		}
		o.Object = &b

	case OBJ_TRANSACTION:

		var cbProbe struct {
			Height *int `json:"height"`
		}

		json.Unmarshal(o.RawObject, &cbProbe)

		if cbProbe.Height != nil {
			var cb CoinbaseTransaction
			if err := json.Unmarshal(o.RawObject, &cb); err != nil {
				return fmt.Errorf("failed to parse coinbase transaction: %w", err)
			}
			o.Object = &cb
		} else {
			var tx Transaction
			if err := json.Unmarshal(o.RawObject, &tx); err != nil {
				return fmt.Errorf("failed to parse transaction: %w", err)
			}
			o.Object = &tx
		}
	default:
		return fmt.Errorf("unknown object type: %s", typeProbe.Type)
	}
	return nil
}
