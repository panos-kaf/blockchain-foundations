package peer

import (
	"marabu/internal/messages"
)

func (p *Peer) handleHello(msg *messages.HelloSchema) {}

func (p *Peer) handleGetPeers() {}

func (p *Peer) handleGetObject(msg *messages.GetObjectSchema) {}

func (p *Peer) handleIHaveObject(msg *messages.IHaveObjectSchema) {}

func (p *Peer) handleObject(msg *messages.ObjectSchema) {}

func (p *Peer) handleGetMempool() {}

func (p *Peer) handleMempool(msg *messages.MempoolSchema) {}

func (p *Peer) handleGetChainTip() {}

func (p *Peer) handleChainTip(msg *messages.ChainTipSchema) {}
