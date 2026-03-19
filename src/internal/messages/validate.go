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

// -- String Validators --

func ValidateStringMaxLen(val string, fieldName string, max int) (error, ErrorCode) {
	if len(val) > max {
		return fmt.Errorf("%s exceeds maximum length of %d (got %d)", fieldName, max, len(val)), E_INVALID_FORMAT
	}
	return nil, E_NONE
}

// JSON schema validation for HashID/Signature checks this already
// func ValidateStringExactLen(val string, fieldName string, exact int) (error, ErrorCode) {
// 	if len(val) != exact {
// 		return fmt.Errorf("%s must be exactly %d characters (got %d)", fieldName, exact, len(val)), E_INVALID_FORMAT
// 	}
// 	return nil, E_NONE
// }

// -- Array/Slice Validators --

func ValidateStringSliceMaxLen(arr []string, fieldName string, max int) (error, ErrorCode) {
	if len(arr) > max {
		return fmt.Errorf("%s exceeds maximum array length of %d (got %d)", fieldName, max, len(arr)), E_INVALID_FORMAT
	}
	return nil, E_NONE
}

// Ensure every string in a slice doesn't exceed a specific length (e.g., studentids max 128)
func ValidateStringSliceElementMaxLen(arr []string, fieldName string, max int) (error, ErrorCode) {
	for i, v := range arr {
		if err, code := ValidateStringMaxLen(v, fmt.Sprintf("%s[%d]", fieldName, i), max); err != nil {
			return err, code
		}
	}
	return nil, E_NONE
}

// -- Number Validators --

func ValidateNonNegativeInt(val int, fieldName string) (error, ErrorCode) {
	if val < 0 {
		return fmt.Errorf("%s must be non-negative (got %d)", fieldName, val), E_INVALID_FORMAT
	}
	return nil, E_NONE
}

// -- Regex Validators --

func ValidateVersionString(val string) (error, ErrorCode) {
	if !versionRegex.MatchString(val) {
		return fmt.Errorf("invalid version format: %s", val), E_INVALID_FORMAT
	}
	return nil, E_NONE
}

func ValidatePeerFormat(val string) (error, ErrorCode) {
	if !peerRegex.MatchString(val) {
		return fmt.Errorf("%s", val), E_INVALID_FORMAT
	}
	return nil, E_NONE
}

func ValidatePeers(peers []string) ([]string, error, ErrorCode) {
	var validPeers []string
	var invalid []string
	for _, peer := range peers {
		peer = strings.TrimSpace(peer)
		err, _ := ValidatePeerFormat(peer)
		if err == nil {
			validPeers = append(validPeers, peer)
		} else {
			invalid = append(invalid, peer)
		}
	}
	if len(invalid) > 0 {
		return validPeers, fmt.Errorf("some peers were invalid and ignored: %v", invalid), E_INVALID_FORMAT
	}
	return validPeers, nil, E_NONE
}

// -- Message Type Validators --

func (h *HelloSchema) Validate() (error, ErrorCode) {
	if err, code := ValidateVersionString(h.Version); err != nil {
		return err, code
	}
	if h.Agent != nil {
		return ValidateStringMaxLen(*h.Agent, "agent", 1000)
	}
	return nil, E_NONE
}

func (e *ErrorSchema) Validate() (error, ErrorCode) {
	return ValidateStringMaxLen(e.Description, "description", 1000)
}

func (g *GetPeersSchema) Validate() (error, ErrorCode) {
	return nil, E_NONE
}

func (p *PeersSchema) Validate() (error, ErrorCode) {
	if err, code := ValidateStringSliceMaxLen(p.Peers, "peers", 1000); err != nil {
		return err, code
	}
	if err, code := ValidateStringSliceElementMaxLen(p.Peers, "peers", 1000); err != nil {
		return err, code
	}
	peers, err, code := ValidatePeers(p.Peers)
	p.Peers = peers

	return err, code
}

func (g *GetObjectSchema) Validate() (error, ErrorCode) {
	return nil, E_NONE
}

func (i *IHaveObjectSchema) Validate() (error, ErrorCode) {
	return nil, E_NONE
}

func (o *ObjectSchema) Validate() (error, ErrorCode) {

	if o.Object == nil {
		return fmt.Errorf("object could not get parsed"), E_INVALID_FORMAT
	}

	return o.Object.Validate()
}

func (g *GetMempoolSchema) Validate() (error, ErrorCode) {
	return nil, E_NONE
}

func (m *MempoolSchema) Validate() (error, ErrorCode) {
	return nil, E_NONE
}

func (g *GetChainTipSchema) Validate() (error, ErrorCode) {
	return nil, E_NONE
}

func (c *ChainTipSchema) Validate() (error, ErrorCode) {
	return nil, E_NONE
}
