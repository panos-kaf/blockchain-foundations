package peer

import (
	"marabu/internal/messages"
)

func (p *Peer) handleHello(msg *messages.HelloSchema) {
	if msg.Agent != nil {
		p.log("Received HELLO from peer: " + p.id + " (agent: " + *msg.Agent + ")" + " version: " + msg.Version)
	} else {
		p.log("Received HELLO from peer: " + p.id + " version: " + msg.Version)
	}
	p.handshakeComplete = true
}

func (p *Peer) handleGetPeers() {
	p.log("Received GETPEERS from peer: " + p.id)
	peers := make([]string, 0, len(knownPeers))
	for peer := range knownPeers {
		peers = append(peers, peer)
	}
	p.SendMessage(messages.MakePeersMessage(peers))
}

func (p *Peer) handlePeers(msg *messages.PeersSchema) {
	p.log("Received PEERS from peer: " + p.id)
	AppendPeers(msg.Peers, p.id)
}

func (p *Peer) handleGetObject(msg *messages.GetObjectSchema) {}

func (p *Peer) handleIHaveObject(msg *messages.IHaveObjectSchema) {}

func (p *Peer) handleObject(msg *messages.ObjectSchema) {}

func (p *Peer) handleGetMempool() {}

func (p *Peer) handleMempool(msg *messages.MempoolSchema) {}

func (p *Peer) handleGetChainTip() {}

func (p *Peer) handleChainTip(msg *messages.ChainTipSchema) {}
