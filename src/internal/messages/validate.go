package messages

import (
	"fmt"
)

// func ValidatePeers(peers []string) ([]string, error, ErrorCode) {
// 	var validPeers []string
// 	var invalid []string
// 	for _, peer := range peers {
// 		peer = strings.TrimSpace(peer)

// 		err, _ := ValidatePeerFormat(peer)
// 		if err == nil {
// 			validPeers = append(validPeers, peer)
// 		} else {
// 			invalid = append(invalid, peer)
// 		}
// 	}
// 	if len(invalid) > 0 {
// 		return validPeers, fmt.Errorf("some peers were invalid and ignored: %v", invalid), E_INVALID_FORMAT
// 	}
// 	return validPeers, nil, E_NONE
// }

// -- Message Type Validators --

func (h *HelloSchema) Validate() (error, ErrorCode) {
	return nil, E_NONE
}

func (e *ErrorSchema) Validate() (error, ErrorCode) {
	return nil, E_NONE
}

func (g *GetPeersSchema) Validate() (error, ErrorCode) {
	return nil, E_NONE
}

func (p *PeersSchema) Validate() (error, ErrorCode) {
	return nil, E_NONE
	// Obsolete (?)
	// peers, err, code := ValidatePeers(p.Peers)
	// p.Peers = peers

	// return err, code
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
