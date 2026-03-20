package messages

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
)

// Reusable compiled regexes for performance
var (
	versionRegex = regexp.MustCompile(`^0\.10\.[0-9]+$`)
	peerRegex    = regexp.MustCompile(`^((?:\d{1,3}\.){3}\d{1,3}|\[[a-fA-F0-9:]+\]|[a-fA-F0-9:]+|[a-zA-Z0-9.-]+):[0-9]{1,5}$`)
)

const maxArrLen = 1000

func ValidateVersionString(val string) (error, ErrorCode) {

	return nil, E_NONE
}

var messageTypeRegistry = map[string]reflect.Type{
	string(MSG_HELLO):       reflect.TypeOf(HelloMessage{}),
	string(MSG_ERROR):       reflect.TypeOf(ErrorMessage{}),
	string(MSG_GETPEERS):    reflect.TypeOf(GetPeersMessage{}),
	string(MSG_PEERS):       reflect.TypeOf(PeersMessage{}),
	string(MSG_GETOBJECT):   reflect.TypeOf(GetObjectMessage{}),
	string(MSG_IHAVEOBJECT): reflect.TypeOf(IHaveObjectMessage{}),
	string(MSG_OBJECT):      reflect.TypeOf(ObjectMessage{}),
	string(MSG_GETMEMPOOL):  reflect.TypeOf(GetMempoolMessage{}),
	string(MSG_MEMPOOL):     reflect.TypeOf(MempoolMessage{}),
	string(MSG_GETCHAINTIP): reflect.TypeOf(GetChainTipMessage{}),
	string(MSG_CHAINTIP):    reflect.TypeOf(ChainTipMessage{}),
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

func (v *T_Version) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("invalid version format: %w", err)
	}
	if !versionRegex.MatchString(s) {
		return fmt.Errorf("invalid version format: %s", s)
	}
	*v = T_Version(s)
	return nil
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

func (p *T_Peer) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("invalid peer format: %w", err)
	}
	if !peerRegex.MatchString(s) {
		return fmt.Errorf("invalid peer format: %s", s)
	}
	*p = T_Peer(s)
	return nil
}

// Custom UnmarshalJSON for Hash IDs to enforce length and hex format
func (h *T_HashID) UnmarshalJSON(data []byte) error {
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

	*h = T_HashID(s)
	return nil
}

// Custom UnmarshalJSON for Signatures to enforce length and hex format
func (sig *T_Signature) UnmarshalJSON(data []byte) error {
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

	*sig = T_Signature(s)
	return nil
}

// Custom UnmarshalJSON for ObjectMessage to handle dynamic inner object types (block, transaction, coinbase transaction)
func (o *ObjectMessage) UnmarshalJSON(data []byte) error {

	type Alias ObjectMessage

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
		var b T_Block
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
			var cb T_CoinbaseTransaction
			if err := json.Unmarshal(o.RawObject, &cb); err != nil {
				return fmt.Errorf("failed to parse coinbase transaction: %w", err)
			}
			o.Object = &cb
		} else {
			var tx T_Transaction
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

func (n *T_BuInt) UnmarshalJSON(data []byte) error {
	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		return fmt.Errorf("invalid integer format: %w", err)
	}

	if i < 0 {
		return fmt.Errorf("integer value must be non-negative, got %d", i)
	}

	*n = T_BuInt(i)
	return nil
}

func (s *T_BuString) UnmarshalJSON(data []byte) error {
	var str string

	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("invalid string format: %w", err)
	}

	if len(str) > maxArrLen {
		return fmt.Errorf("string exceeds maximum length of %d characters, got %d", maxArrLen, len(str))
	}
	*s = T_BuString(str)
	return nil
}

// Generic helper for unmarshaling arrays with length checks
func UnmarshalArray[T any](data []byte, target *[]T, maxLen int, fieldName string) error {
	var arr []T
	if err := json.Unmarshal(data, &arr); err != nil {
		return fmt.Errorf("invalid %s array format: %w", fieldName, err)
	}
	if len(arr) > maxLen {
		return fmt.Errorf("%s array exceeds maximum length of %d elements, got %d", fieldName, maxLen, len(arr))
	}
	*target = arr
	return nil
}

func (arr *T_BuInts) UnmarshalJSON(data []byte) error {
	return UnmarshalArray(data, (*[]T_BuInt)(arr), maxArrLen, "T_BuInts")
}
func (arr *T_BuStrings) UnmarshalJSON(data []byte) error {
	return UnmarshalArray(data, (*[]T_BuString)(arr), maxArrLen, "T_BuStrings")
}
func (arr *T_HashIDs) UnmarshalJSON(data []byte) error {
	return UnmarshalArray(data, (*[]T_HashID)(arr), maxArrLen, "T_HashIDs")
}
func (arr *T_Peers) UnmarshalJSON(data []byte) error {
	return UnmarshalArray(data, (*[]T_Peer)(arr), maxArrLen, "T_Peers")
}
