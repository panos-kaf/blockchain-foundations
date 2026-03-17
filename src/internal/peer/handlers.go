package peer

import (
	"marabu/internal/messages"
	"strconv"
)

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
	peers := make([]string, 0, len(knownPeers))
	for peer := range knownPeers {
		peers = append(peers, peer)
	}
	err := p.SendPeers(peers)
	if err != nil {
		p.logErr(messages.PEERS, err.Error())
	}
}

func (p *Peer) handlePeers(msg *messages.PeersSchema) {
	p.log(messages.PEERS, "Peer "+p.addr+" sent "+strconv.Itoa(len(msg.Peers))+" peers")
	AppendPeers(msg.Peers, p.addr)
}

func (p *Peer) handleGetObject(msg *messages.GetObjectSchema) {

	Log := func(m string) {
		p.log(messages.GETOBJECT, m)
	}
	Err := func(m string) {
		p.logErr(messages.GETOBJECT, m)
	}
	ID := msg.ObjectID
	Log("Peer: " + p.addr + " requested object: " + string(ID))

	exists, err := p.objectManager.Exists(ID)
	if err != nil {
		Err("Error checking if object exists: " + err.Error())
		return
	}
	if exists {
		Log("We have object " + string(ID) + ", sending it to peer " + p.addr)
		obj, err := p.objectManager.Get(ID)
		if err != nil {
			Err("Error retrieving object: " + err.Error())
			return
		}
		err = p.SendObject(ID, obj)
		if err != nil {
			Err("Error sending object: " + err.Error())
		}
	} else {
		Log("We do not have object " + string(ID) + ", cannot fulfill GETOBJECT request from peer " + p.addr)
		p.SendError(messages.UNKNOWN_OBJECT, "Object not found: "+string(ID))
	}
}

func (p *Peer) handleIHaveObject(msg *messages.IHaveObjectSchema) {

	Log := func(m string) {
		p.log(messages.IHAVEOBJECT, m)
	}
	Err := func(m string) {
		p.logErr(messages.IHAVEOBJECT, m)
	}

	ID := msg.ObjectID
	Log("Peer: " + p.addr + "  has object with ID: " + string(ID))

	exists, err := p.objectManager.Exists(ID)
	if err != nil {
		Err("Error checking if object exists: " + err.Error())
		return
	}
	if exists {
		Log("We already have object " + string(ID))
	} else {
		Log("We do not have object " + string(ID) + ", requesting it from peer " + p.addr)
		err := p.SendGetObject(ID)
		if err != nil {
			Err("Error sending GETOBJECT: " + err.Error())
		}
	}
}

func (p *Peer) handleObject(msg *messages.ObjectSchema) {

	Log := func(m string) {
		p.log(messages.OBJECT, m)
	}
	Err := func(m string) {
		p.logErr(messages.OBJECT, m)
	}

	errorCode, err := p.ValidateObject(msg.Object, msg.ObjectID)
	if err != nil {
		Err("Received invalid object from peer " + p.addr + ": " + err.Error())
		p.SendError(errorCode, "Invalid object: "+err.Error())
		return
	}

	ID := msg.ObjectID
	IDStr := string(ID)

	Log("Received OBJECT with ID " + IDStr + " from peer: " + p.addr)

	exists, err := p.objectManager.Exists(ID)
	if err != nil {
		Err("Error checking if object exists: " + err.Error())
		return
	}
	if exists {
		Log("We already have object " + IDStr + ", ignoring received object.")
	} else {
		Log("Storing new object with ID " + IDStr)

		_, err := p.objectManager.Put(msg.Object)
		if err != nil {
			Err("Error storing object: " + err.Error())
			return
		}
		Log("Object stored successfully with ID " + IDStr)

		// gossip!
		advertisement, err := messages.MakeIHaveObjectMessage(ID)
		if err != nil {
			Err("Error creating IHAVEOBJECT message: " + err.Error())
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
