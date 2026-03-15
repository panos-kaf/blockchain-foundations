package messages

import (
	"fmt"
	"regexp"
	"strings"
)

// Reusable compiled regexes for performance
var (
	versionRegex = regexp.MustCompile(`^0\.10\.[0-9]+$`)
	peerRegex    = regexp.MustCompile(`^((?:\d{1,3}\.){3}\d{1,3}|\[[a-fA-F0-9:]+\]|[a-fA-F0-9:]+|[a-zA-Z0-9.-]+):[0-9]{1,5}$`)
)

// -- Type Validators --

func ValidateMessageType(val MessageType) error {
	switch val {
	case HELLO, ERROR, GETPEERS, PEERS, GETOBJECT, IHAVEOBJECT, OBJECT, GETMEMPOOL, GETCHAINTIP, CHAINTIP:
		return nil
	default:
		return fmt.Errorf("invalid message type: %s", val)
	}
}

func ValidateErrorCode(val ErrorCode) error {
	switch val {
	case INTERNAL_ERROR, INVALID_FORMAT, UNKNOWN_OBJECT, UNFINDABLE_OBJECT, INVALID_HANDSHAKE, INVALID_TX_OUTPOINT, INVALID_TX_SIGNATURE, INVALID_TX_CONSERVATION, INVALID_BLOCK_COINBASE, INVALID_BLOCK_TIMESTAMP, INVALID_BLOCK_POW, INVALID_GENESIS:
		return nil
	default:
		return fmt.Errorf("invalid error code: %s", val)
	}
}

func ValidateObjectType(val ObjectType) error {
	switch val {
	case TRANSACTION, BLOCK:
		return nil
	default:
		return fmt.Errorf("invalid object type: %s", val)
	}
}

// -- String Validators --

func ValidateStringMaxLen(val string, fieldName string, max int) error {
	if len(val) > max {
		return fmt.Errorf("%s exceeds maximum length of %d (got %d)", fieldName, max, len(val))
	}
	return nil
}

func ValidateStringExactLen(val string, fieldName string, exact int) error {
	if len(val) != exact {
		return fmt.Errorf("%s must be exactly %d characters (got %d)", fieldName, exact, len(val))
	}
	return nil
}

// -- Array/Slice Validators --

func ValidateStringSliceMaxLen(arr []string, fieldName string, max int) error {
	if len(arr) > max {
		return fmt.Errorf("%s exceeds maximum array length of %d (got %d)", fieldName, max, len(arr))
	}
	return nil
}

// Ensure every string in a slice doesn't exceed a specific length (e.g., studentids max 128)
func ValidateStringSliceElementMaxLen(arr []string, fieldName string, max int) error {
	for i, v := range arr {
		if err := ValidateStringMaxLen(v, fmt.Sprintf("%s[%d]", fieldName, i), max); err != nil {
			return err
		}
	}
	return nil
}

// -- Number Validators --

func ValidateNonNegativeInt(val int, fieldName string) error {
	if val < 0 {
		return fmt.Errorf("%s must be non-negative (got %d)", fieldName, val)
	}
	return nil
}

// -- Regex Validators --

func ValidateVersionString(val string) error {
	if !versionRegex.MatchString(val) {
		return fmt.Errorf("invalid version format: %s", val)
	}
	return nil
}

func ValidatePeerFormat(val string) error {
	if !peerRegex.MatchString(val) {
		return fmt.Errorf("%s", val)
	}
	return nil
}

func ValidatePeers(peers []string) ([]string, error) {
	var validPeers []string
	var invalid []string
	for _, peer := range peers {
		peer = strings.TrimSpace(peer)
		if ValidatePeerFormat(peer) == nil {
			validPeers = append(validPeers, peer)
		} else {
			invalid = append(invalid, peer)
		}
	}
	if len(invalid) > 0 {
		return validPeers, fmt.Errorf("some peers were invalid and ignored: %v", invalid)
	}
	return validPeers, nil
}

// -- Message Type Validators --

func (h *HelloSchema) Validate() error {

	if err := ValidateMessageType(h.Type); err != nil {
		return err
	}
	if err := ValidateVersionString(h.Version); err != nil {
		return err
	}
	if h.Agent != nil {
		return ValidateStringMaxLen(*h.Agent, "agent", 1000)
	}
	return nil
}

func (e *ErrorSchema) Validate() error {

	if err := ValidateMessageType(e.Type); err != nil {
		return err
	}
	if err := ValidateErrorCode(e.Name); err != nil {
		return err
	}
	return ValidateStringMaxLen(e.Description, "description", 1000)
}

func (g *GetPeersSchema) Validate() error {
	return ValidateMessageType(g.Type)
}

func (p *PeersSchema) Validate() error {
	if err := ValidateMessageType(p.Type); err != nil {
		return err
	}
	if err := ValidateStringSliceMaxLen(p.Peers, "peers", 1000); err != nil {
		return err
	}
	if err := ValidateStringSliceElementMaxLen(p.Peers, "peers", 1000); err != nil {
		return err
	}
	peers, err := ValidatePeers(p.Peers)
	p.Peers = peers

	return err
}

func (g *GetObjectSchema) Validate() error {
	return ValidateMessageType(g.Type)
}

func (i *IHaveObjectSchema) Validate() error {
	return ValidateMessageType(i.Type)
}

func (o *ObjectSchema) Validate() error {
	if err := ValidateMessageType(o.MessageType()); err != nil {
		return err
	}

	if o.Object == nil {
		return fmt.Errorf("object could not get parsed")
	}

	ID, err := HashObject(o.Object)
	if err != nil {
		return err
	}

	if ID != o.ObjectID {
		return fmt.Errorf("object ID does not match hash of object content")
	}

	return o.Object.Validate()
}

func (g *GetMempoolSchema) Validate() error {
	return ValidateMessageType(g.Type)
}

func (m *MempoolSchema) Validate() error {
	return ValidateMessageType(m.Type)
	// return ValidateStringSliceMaxLen(m.Txids, "txids", 1000)
}

func (g *GetChainTipSchema) Validate() error {
	return ValidateMessageType(g.Type)
}

func (c *ChainTipSchema) Validate() error {
	return ValidateMessageType(c.Type)
}
