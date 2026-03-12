package messages

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var messageTypeRegistry = map[string]reflect.Type{
	string(HELLO):       reflect.TypeOf(HelloSchema{}),
	string(ERROR):       reflect.TypeOf(ErrorSchema{}),
	string(GETPEERS):    reflect.TypeOf(GetPeersSchema{}),
	string(PEERS):       reflect.TypeOf(PeersSchema{}),
	string(GETOBJECT):   reflect.TypeOf(GetObjectSchema{}),
	string(IHAVEOBJECT): reflect.TypeOf(IHaveObjectSchema{}),
	string(OBJECT):      reflect.TypeOf(ObjectSchema{}),
	string(GETMEMPOOL):  reflect.TypeOf(GetMempoolSchema{}),
	string(MEMPOOL):     reflect.TypeOf(MempoolSchema{}),
	string(GETCHAINTIP): reflect.TypeOf(GetChainTipSchema{}),
	string(CHAINTIP):    reflect.TypeOf(ChainTipSchema{}),
}

func UnmarshalMessage(raw string) (Message, error) {
	var probe map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &probe); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}
	typeVal, ok := probe["type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'type' field in message")
	}

	typ, found := messageTypeRegistry[typeVal]
	if !found {
		return nil, fmt.Errorf("unknown message type: %s", typeVal)
	}

	msgPtr := reflect.New(typ)
	if err := json.Unmarshal([]byte(raw), msgPtr.Interface()); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s message: %w", typeVal, err)
	}

	// Return as Message interface
	if m, ok := msgPtr.Interface().(Message); ok {
		return m, nil
	}
	return nil, fmt.Errorf("type %s does not implement Message interface", typeVal)
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

	if len(s) != 64 {
		return fmt.Errorf("invalid signature: must be exactly 64 characters, got %d", len(s))
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
	case BLOCK:
		var b Block
		if err := json.Unmarshal(o.RawObject, &b); err != nil {
			return fmt.Errorf("failed to unmarshal block object: %w", err)
		}
		o.Object = &b

	case TRANSACTION:

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
