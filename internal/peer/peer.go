package peer

import (
	"bufio"
	"marabu/internal/messages"
	"net"
	"sync"
)

var (
	connectedPeers      = make(map[string]*Peer)
	connectedPeersMutex sync.Mutex
)

type Peer struct {
	conn              net.Conn
	id                string
	buffer            []byte
	handshakeComplete bool
	onDisconnect      func()
	onLog             func(string)
	onLogErr          func(string)
	role              string
}

func NewPeer(conn net.Conn, role string, onDisconnect func(), onLog func(string), onLogErr func(string)) *Peer {
	id := conn.RemoteAddr().String()
	p := &Peer{
		conn:         conn,
		id:           id,
		buffer:       make([]byte, 0),
		onLog:        onLog,
		onLogErr:     onLogErr,
		onDisconnect: onDisconnect,
		role:         role,
	}

	connectedPeersMutex.Lock()
	connectedPeers[id] = p
	connectedPeersMutex.Unlock()

	go p.initializeSocket()

	return p
}

func (p *Peer) initializeSocket() {
	reader := bufio.NewReader(p.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			p.logErr("Disconnected: " + err.Error())
			connectedPeersMutex.Lock()
			delete(connectedPeers, p.id)
			connectedPeersMutex.Unlock()
			if p.onDisconnect != nil {
				p.onDisconnect()
			}
			return
		}
		p.handleMessage(line)
	}
}

func (p *Peer) handleMessage(raw string) {
	// Unmarshal and validate message
	var msg messages.Message
	msg, err := messages.UnmarshalMessage(raw)
	if err != nil {
		p.logErr("Invalid message: " + err.Error())
		p.SendError(messages.INVALID_FORMAT, "Could not parse message as JSON")
		if !p.handshakeComplete {
			p.conn.Close()
		}
		return
	}

	// Dispatch based on type
	switch m := msg.(type) {
	case *messages.HelloSchema:
		p.handleHello(m)
	case *messages.GetPeersSchema:
		p.handleGetPeers()
	case *messages.GetObjectSchema:
		p.handleGetObject(m)
	case *messages.IHaveObjectSchema:
		p.handleIHaveObject(m)
	case *messages.ObjectSchema:
		p.handleObject(m)
	case *messages.GetMempoolSchema:
		p.handleGetMempool()
	case *messages.MempoolSchema:
		p.handleMempool(m)
	case *messages.GetChainTipSchema:
		p.handleGetChainTip()
	case *messages.ChainTipSchema:
		p.handleChainTip(m)
	default:
		p.logErr("Unknown message type")
		p.SendError(messages.INVALID_FORMAT, "Unknown protocol message")
		p.conn.Close()
	}
}

func (p *Peer) log(msg string) {
	if p.onLog != nil {
		p.onLog(msg)
	} else {
		println("[" + p.role + ":" + p.id + "] " + msg)
	}
}
func (p *Peer) logErr(msg string) {
	if p.onLogErr != nil {
		p.onLogErr(msg)
	} else {
		println("[" + p.role + ":" + p.id + "] ERROR: " + msg)
	}
}

func Broadcast(msg string) {
	connectedPeersMutex.Lock()
	defer connectedPeersMutex.Unlock()
	for _, peer := range connectedPeers {
		peer.SendMessage(msg)
	}
}

func (p *Peer) SendMessage(msg string) error {
	_, err := p.conn.Write([]byte(msg))
	return err
}

func (p *Peer) SendError(name messages.ErrorCode, description string) error {
	msg, err := messages.MakeErrorMessage(name, description)
	if err != nil {
		return err
	}
	return p.SendMessage(msg)
}
