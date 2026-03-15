package peer

import (
	"marabu/internal/messages"
	"strconv"
)

type HashID = messages.HashID

func (p *Peer) handleHello(msg *messages.HelloSchema) {
	if msg.Agent != nil {
		p.log(msg.Type, *msg.Agent+" ("+p.addr+") says hello, version: "+msg.Version)
	} else {
		p.log(msg.Type, "Peer "+p.addr+" says hello, version: "+msg.Version)
	}
	p.handshakeComplete = true
}

func (p *Peer) handleError(msg *messages.ErrorSchema) {
	p.log(msg.Type, string(msg.Name)+", peer: "+p.addr+", description: "+msg.Description+")")
}

func (p *Peer) handleGetPeers() {
	p.log(messages.GETPEERS, "Peer "+p.addr+" requested peers")
	peers := make([]string, 0, len(knownPeers))
	for peer := range knownPeers {
		peers = append(peers, peer)
	}
	err := p.SendPeers(peers)
	if err != nil {
		p.logErr(err.Error())
	}
}

func (p *Peer) handlePeers(msg *messages.PeersSchema) {
	p.log(messages.PEERS, "Peer "+p.addr+" sent "+strconv.Itoa(len(msg.Peers))+" peers")
	AppendPeers(msg.Peers, p.addr)
}

func (p *Peer) handleGetObject(msg *messages.GetObjectSchema) {
	ID := msg.ID
	p.log(messages.GETOBJECT, "Peer: "+p.addr+" requested object: "+string(ID))

	exists, err := p.objectManager.Exists(ID)
	if err != nil {
		p.logErr("Error checking if object exists: " + err.Error())
		return
	}
	if exists {
		p.log(messages.GETOBJECT, "We have object "+string(ID)+", sending it to peer "+p.addr)
		obj, err := p.objectManager.Get(ID)
		if err != nil {
			p.logErr("Error retrieving object: " + err.Error())
			return
		}
		err = p.SendObject(ID, obj)
		if err != nil {
			p.logErr("Error sending object: " + err.Error())
		}
	} else {
		p.log(messages.GETOBJECT, "We do not have object "+string(ID)+", cannot fulfill GETOBJECT request from peer "+p.addr)
		p.SendError(messages.UNKNOWN_OBJECT, "Object not found: "+string(ID))
	}
}

func (p *Peer) handleIHaveObject(msg *messages.IHaveObjectSchema) {

	ID := msg.ID
	p.log(messages.IHAVEOBJECT, "Peer: "+p.addr+"  has object with ID: "+string(ID))

	exists, err := p.objectManager.Exists(ID)
	if err != nil {
		p.logErr("Error checking if object exists: " + err.Error())
		return
	}
	if exists {
		p.log(messages.IHAVEOBJECT, "We already have object "+string(ID))
	} else {
		p.log(messages.IHAVEOBJECT, "We do not have object "+string(ID)+", requesting it from peer "+p.addr)
		err := p.SendGetObject(ID)
		if err != nil {
			p.logErr("Error sending GETOBJECT: " + err.Error())
		}
	}
}

func (p *Peer) handleObject(msg *messages.ObjectSchema) {

	ID := msg.ObjectID
	p.log(messages.OBJECT, "Received OBJECT with ID "+string(ID)+" from peer: "+p.addr)

	exists, err := p.objectManager.Exists(ID)
	if err != nil {
		p.logErr("Error checking if object exists: " + err.Error())
		return
	}
	if exists {
		p.log(messages.OBJECT, "We already have object "+string(ID)+", ignoring received object.")
	} else {
		p.log(messages.OBJECT, "Storing new object with ID "+string(ID))
		_, err := p.objectManager.Put(msg.Object)
		if err != nil {
			p.logErr("Error storing object: " + err.Error())
			return
		}
		p.log(messages.OBJECT, "Object stored successfully with ID "+string(ID))

		// gossip!
		advertisement, err := messages.MakeIHaveObjectMessage(ID)
		if err != nil {
			p.logErr("Error creating IHAVEOBJECT message: " + err.Error())
			return
		}
		Broadcast(messages.IHAVEOBJECT, advertisement, err)
	}
}

func (p *Peer) handleGetMempool() {
	p.log(messages.GETMEMPOOL, "not handled yet")
}

func (p *Peer) handleMempool(msg *messages.MempoolSchema) {
	p.log(msg.Type, "not handled yet")
}

func (p *Peer) handleGetChainTip() {
	p.log(messages.GETCHAINTIP, "not handled yet")
}

func (p *Peer) handleChainTip(msg *messages.ChainTipSchema) {
	p.log(msg.Type, "not handled yet")
}
