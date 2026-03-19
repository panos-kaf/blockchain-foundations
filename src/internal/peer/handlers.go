package peer

import (
	"marabu/internal/crypto"
	"marabu/internal/messages"
	"strconv"
)

func (p *Peer) handleHello(msg *HelloSchema) {
	if msg.Agent != nil {
		p.log(msg.Type, *msg.Agent+" ("+p.addr+") says hello, version: "+msg.Version)
	} else {
		p.log(msg.Type, "Peer "+p.addr+" says hello, version: "+msg.Version)
	}
	p.handshakeComplete = true
}

func (p *Peer) handleError(msg *ErrorSchema) {
	p.log(msg.Type, string(msg.Name)+", peer: "+p.addr+", description: "+msg.Description+")")
}

func (p *Peer) handleGetPeers() {
	peers := make([]string, 0, len(knownPeers))
	for peer := range knownPeers {
		peers = append(peers, peer)
	}
	err := p.SendPeers(peers)
	if err != nil {
		p.logErr(MSG_PEERS, err.Error())
	}
}

func (p *Peer) handlePeers(msg *PeersSchema) {
	p.log(MSG_PEERS, "Peer "+p.addr+" sent "+strconv.Itoa(len(msg.Peers))+" peers")
	AppendPeers(msg.Peers, p.addr)
}

func (p *Peer) handleGetObject(msg *GetObjectSchema) {

	Log := func(m string) {
		p.log(MSG_GETOBJECT, m)
	}
	Err := func(m string) {
		p.logErr(MSG_GETOBJECT, m)
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
		err = p.SendObject(obj)
		if err != nil {
			Err("Error sending object: " + err.Error())
		}
	} else {
		Log("We do not have object " + string(ID) + ", cannot fulfill MSG_GETOBJECT request from peer " + p.addr)
		p.SendError(E_UNKNOWN_OBJECT, "Object not found: "+string(ID))
	}
}

func (p *Peer) handleIHaveObject(msg *IHaveObjectSchema) {

	Log := func(m string) {
		p.log(MSG_IHAVEOBJECT, m)
	}
	Err := func(m string) {
		p.logErr(MSG_IHAVEOBJECT, m)
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
			Err("Error sending MSG_GETOBJECT: " + err.Error())
		}
	}
}

func (p *Peer) handleObject(msg *ObjectSchema) {

	Log := func(m string) {
		p.log(MSG_OBJECT, m)
	}
	Err := func(m string) {
		p.logErr(MSG_OBJECT, m)
	}

	errorCode, err := p.ValidateObject(msg.Object)
	if err != nil {
		Err("Received invalid object from peer " + p.addr + ": " + err.Error())
		p.SendError(errorCode, "Invalid object: "+err.Error())
		return
	}

	ID, err := crypto.HashObject(msg.Object)
	if err != nil {
		Err("Error hashing object: " + err.Error())
		return
	}

	hashID := HashID(ID)

	Log("Received MSG_OBJECT with ID " + ID + " from peer: " + p.addr)

	exists, err := p.objectManager.Exists(hashID)
	if err != nil {
		Err("Error checking if object exists: " + err.Error())
		return
	}
	if exists {
		Log("We already have object " + ID + ", ignoring received object.")
	} else {
		Log("Storing new object with ID " + ID)

		_, err := p.objectManager.Put(msg.Object)
		if err != nil {
			Err("Error storing object: " + err.Error())
			return
		}
		Log("Object stored successfully with ID " + ID)

		// gossip!
		advertisement, err := messages.MakeIHaveObjectMessage(hashID)
		if err != nil {
			Err("Error creating MSG_IHAVEOBJECT message: " + err.Error())
			return
		}
		Broadcast(MSG_IHAVEOBJECT, E_NONE, advertisement, err)
	}
}

func (p *Peer) handleGetMempool() {
	p.log(MSG_GETMEMPOOL, "not handled yet")
}

func (p *Peer) handleMempool(msg *MempoolSchema) {
	p.log(msg.Type, "not handled yet")
}

func (p *Peer) handleGetChainTip() {
	p.log(MSG_GETCHAINTIP, "not handled yet")
}

func (p *Peer) handleChainTip(msg *ChainTipSchema) {
	p.log(msg.Type, "not handled yet")
}
