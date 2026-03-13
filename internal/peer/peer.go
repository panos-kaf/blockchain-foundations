package peer

import (
	"bufio"
	"fmt"
	"marabu/internal/messages"
	"net"
	"strconv"
	"strings"
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

	if len(strings.TrimSpace(raw)) == 0 {
		p.log("Received empty message")
		return
	}

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

	if err := msg.Validate(); err != nil {
		p.logErr("Message validation failed: " + err.Error())
		p.SendError(messages.INVALID_FORMAT, "Message validation failed: "+err.Error())

		if !p.handshakeComplete {
			p.conn.Close()
		}
		return
	}

	if !p.handshakeComplete && msg.MessageType() != messages.HELLO {
		p.logErr("Expected HELLO message first")
		p.SendError(messages.INVALID_HANDSHAKE, "Handshake not completed, expected hello message but received "+string(msg.MessageType()))
		p.conn.Close()
		return
	}

	// Dispatch based on type
	switch m := msg.(type) {
	case *messages.HelloSchema:
		p.handleHello(m)
	case *messages.GetPeersSchema:
		p.handleGetPeers()
	case *messages.PeersSchema:
		p.handlePeers(m)
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
		fmt.Println("[" + p.role + ":" + p.id + "] " + msg)
	}
}

func (p *Peer) logErr(msg string) {
	if p.onLogErr != nil {
		p.onLogErr(msg)
	} else {
		fmt.Println("[" + p.role + ":" + p.id + "] ERROR: " + msg)
	}
}

func (p *Peer) SendMessage(msg string, mkErr error) error {
	if mkErr != nil {
		return mkErr
	}
	_, err := p.conn.Write([]byte(msg))
	return err
}

func (p *Peer) SendError(name messages.ErrorCode, description string) error {
	msg, err := messages.MakeErrorMessage(name, description)

	return p.SendMessage(msg, err)
}

func (p *Peer) Greet() {
	p.SendMessage(messages.MakeHelloMessage())
	p.SendMessage(messages.MakeGetPeersMessage())
}

func Broadcast(msg string, mkErr error) {
	connectedPeersMutex.Lock()
	defer connectedPeersMutex.Unlock()
	for _, peer := range connectedPeers {
		peer.SendMessage(msg, mkErr)
	}
}

func StartServer(port int) error {

	addr := net.JoinHostPort("", strconv.Itoa(port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	fmt.Printf("Peer server listening on port %d...\n", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %s\n", err)
			continue
		}
		addr := conn.RemoteAddr().String()
		fmt.Printf("Accepted connection from %s\n", addr)

		onLog := func(msg string) {
			fmt.Printf("Peer %s: %s\n", addr, msg)
		}
		onLogErr := func(errMsg string) {
			fmt.Printf("Peer %s error: %s\n", addr, errMsg)
		}
		onDisconnect := func() {
			fmt.Printf("Peer %s disconnected\n", addr)
		}

		p := NewPeer(conn, "server", onDisconnect, onLog, onLogErr)

		p.Greet()
	}
}

func StartClient(host string, port int, id int, onClose func()) error {

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		if onClose != nil {
			onClose()
		}
		return err
	}

	onLog := func(msg string) {
		fmt.Printf("Peer client %d: %s\n", id, msg)
	}
	onLogErr := func(errMsg string) {
		fmt.Printf("Peer client %d error: %s\n", id, errMsg)
	}

	onDisconnect := func() {
		fmt.Printf("Peer client %d disconnected\n", id)
	}

	onLog(fmt.Sprintf("Peer client %d connected to server at %s:%d\n", id, host, port))

	p := NewPeer(conn, fmt.Sprintf("client-%d", id), onDisconnect, onLog, onLogErr)

	p.Greet()

	return nil
}
