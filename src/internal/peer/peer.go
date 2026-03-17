package peer

import (
	"bufio"
	"fmt"
	"io"
	"marabu/internal/logs"
	"marabu/internal/messages"
	"marabu/internal/object"
	"net"
	"strconv"
	"strings"
	"sync"
)

var (
	connectedPeers      = make(map[string]*Peer)
	connectedPeersMutex sync.Mutex
	connectedPeersCnt   = 0
)

type Peer struct {
	conn              net.Conn
	addr              string
	ID                int
	buffer            []byte
	handshakeComplete bool
	onDisconnect      func()
	onLog             func(messages.MessageType, string)
	onLogErr          func(messages.MessageType, string)
	onLogMessage      func(messages.MessageType, bool)
	role              string
	objectManager     *object.ObjectManager
}

// NewPeer creates a new Peer instance for a given network connection.
// It initializes the peer's state and starts a goroutine
// to handle incoming messages from the connection.
func NewPeer(conn net.Conn,
	role string,
	objectManager *object.ObjectManager,
	onDisconnect func(),
	onLog func(messages.MessageType, string),
	onLogErr func(messages.MessageType, string),
	onLogMessage func(messages.MessageType, bool)) *Peer {

	addr := conn.RemoteAddr().String()
	p := &Peer{
		conn:          conn,
		addr:          addr,
		buffer:        make([]byte, 0),
		onLog:         onLog,
		onLogErr:      onLogErr,
		onLogMessage:  onLogMessage,
		onDisconnect:  onDisconnect,
		role:          role,
		objectManager: objectManager,
	}

	connectedPeersMutex.Lock()
	connectedPeers[addr] = p
	p.ID = connectedPeersCnt
	connectedPeersCnt++
	connectedPeersMutex.Unlock()

	go p.initializeSocket()

	return p
}

// initializeSocket starts a goroutine to read messages from the peer's connection.
// It continuously reads lines from the connection, and for each line, it calls handleMessage.
// On error it disconnects and removes the peer from the connectedPeers map.
func (p *Peer) initializeSocket() {
	reader := bufio.NewReader(p.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {

			connectedPeersMutex.Lock()
			delete(connectedPeers, p.addr)
			connectedPeersMutex.Unlock()

			if err != io.EOF {
				p.logErr("", "Disconnected: "+err.Error())
				return
			}
			if p.onDisconnect != nil {
				p.onDisconnect()
				return
			}
		}
		p.handleMessage(line)
	}
}

// handleMessage processes incoming messages,
// ensuring they are valid and dispatching them
// to the appropriate handler based on their type.
func (p *Peer) handleMessage(raw string) {

	if len(strings.TrimSpace(raw)) == 0 {
		p.log("", "Received empty message")
		return
	}

	// Unmarshal and validate message
	var msg messages.Message
	msg, err := messages.UnmarshalMessage(raw)

	if err != nil {
		p.logErr("", "Invalid message: "+err.Error())
		p.SendError(messages.INVALID_FORMAT, "Could not parse message as JSON: "+err.Error())
		if !p.handshakeComplete {
			p.conn.Close()
		}
		return
	}

	if err := msg.Validate(); err != nil {
		p.logErr("", "Message validation failed: "+err.Error())
		p.SendError(messages.INVALID_FORMAT, "Message validation failed: "+err.Error())

		if !p.handshakeComplete {
			p.conn.Close()
		}
		return
	}

	p.logMessage(msg.MessageType(), false)

	if !p.handshakeComplete && msg.MessageType() != messages.HELLO {
		p.logErr("", "Expected HELLO message first")
		p.SendError(messages.INVALID_HANDSHAKE, "Handshake not completed, expected hello message but received "+string(msg.MessageType()))
		p.conn.Close()
		return
	}

	// Dispatch based on type
	switch m := msg.(type) {
	case *messages.HelloSchema:
		p.handleHello(m)
	case *messages.ErrorSchema:
		p.handleError(m)
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
		p.logErr("", "Unknown message type")
		p.SendError(messages.INVALID_FORMAT, "Unknown protocol message")
		p.conn.Close()
	}
}

func globalLog(msg string) {
	logs.GlobalLog(msg)
}

func globalError(msg string) {
	logs.GlobalError(msg)
}

func (p *Peer) log(mtype messages.MessageType, msg string) {
	if p.onLog != nil {
		p.onLog(mtype, msg)
	} else {
		fmt.Println("[" + p.role + ":" + p.addr + "] " + msg)
	}
}

func (p *Peer) logErr(mtype messages.MessageType, msg string) {
	if p.onLogErr != nil {
		p.onLogErr(mtype, msg)
	} else {
		fmt.Println("[" + p.role + ":" + p.addr + "] ERROR: " + msg)
	}
}

func (p *Peer) logMessage(mtype messages.MessageType, sends bool) {
	if p.onLogMessage != nil {
		p.onLogMessage(mtype, sends)
	} else {
		direction := "received"
		if sends {
			direction = "sent"
		}
		fmt.Printf("[%s:%s] %s message: %s\n", p.role, p.addr, direction, mtype)
	}
}

func StartServer(port int, objectManager *object.ObjectManager) error {

	addr := net.JoinHostPort("", strconv.Itoa(port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	logs.GlobalLog(fmt.Sprintf("Peer server listening on port %d...", port))
	for {
		conn, err := ln.Accept()
		if err != nil {
			logs.GlobalError(fmt.Sprintf("Server failed to accept connection: %s", err))
			continue
		}

		addr := conn.RemoteAddr().String()

		p := NewPeer(conn, "server", objectManager, nil, nil, nil, nil)

		p.onLog = func(mtype messages.MessageType, msg string) { logs.ServerLog(mtype, msg, p.ID) }
		p.onLogErr = func(mtype messages.MessageType, msg string) { logs.ServerError(mtype, msg, p.ID) }
		p.onLogMessage = func(mtype messages.MessageType, sends bool) { logs.ServerMessage(mtype, sends, p.ID, p.addr) }
		p.onDisconnect = func() { logs.ServerLog("", fmt.Sprintf("Client at %s disconnected", p.addr), p.ID) }

		p.onLog(messages.HELLO, fmt.Sprintf("Accepted connection from %s", addr))

		p.Greet()
	}
}

func StartClient(host string, port int, objectManager *object.ObjectManager, onClose func()) error {

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		// if onClose != nil {
		// 	onClose()
		// }
		return err
	}

	p := NewPeer(conn, "client", objectManager, nil, nil, nil, nil)

	p.onLog = func(mtype messages.MessageType, msg string) { logs.ClientLog(mtype, msg, p.ID) }
	p.onLogErr = func(mtype messages.MessageType, msg string) { logs.ClientError(mtype, msg, p.ID) }
	p.onLogMessage = func(mtype messages.MessageType, sends bool) { logs.ClientMessage(mtype, sends, p.ID, p.addr) }
	p.onDisconnect = func() { logs.ClientLog("", fmt.Sprintf("Disconnected from server at %s", p.addr), p.ID) }

	p.onLog("", fmt.Sprintf("Connected to server at %s:%d", host, port))

	p.Greet()

	return nil
}
