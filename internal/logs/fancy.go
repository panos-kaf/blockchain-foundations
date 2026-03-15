package logs

import (
	"fmt"
	"log"
	"marabu/internal/messages"
)

const (
	RED     = "\033[31m"
	GREEN   = "\033[32m"
	YELLOW  = "\033[33m"
	BLUE    = "\033[34m"
	MAGENTA = "\033[35m"
	CYAN    = "\033[36m"

	BOLD   = "\033[1m"
	ITALIC = "\033[3m"

	RESET = "\033[0m"
)

func MessageTypeColor(mtype messages.MessageType) string {
	switch mtype {
	case messages.HELLO:
		return GREEN
	case messages.ERROR:
		return RED
	case messages.GETPEERS, messages.PEERS:
		return CYAN
	case messages.GETOBJECT, messages.IHAVEOBJECT, messages.OBJECT:
		return YELLOW
	case messages.GETMEMPOOL, messages.MEMPOOL:
		return BLUE
	case messages.GETCHAINTIP, messages.CHAINTIP:
		return MAGENTA
	default:
		return RESET
	}
}

func ClientLog(mtype messages.MessageType, msg string, id int) {
	if mtype != "" {
		msgcolor := MessageTypeColor(mtype)
		msg = fmt.Sprintf("%s%s[Peer %d]%s[%s] %s%s", BOLD, BLUE, id, msgcolor, mtype, RESET, msg)
	} else {
		msg = fmt.Sprintf("%s%s[Peer %d] %s%s", BOLD, BLUE, id, RESET, msg)
	}

	log.Printf("%s\n", msg)
	// fmt.Printf("%s\n", msg)
}

func ClientError(msg string, id int) {
	ClientLog(messages.ERROR, fmt.Sprintf("%sError: %s%s", RED, msg, RESET), id)
}

func ServerLog(mtype messages.MessageType, msg string, id int) {
	if mtype != "" {
		msgcolor := MessageTypeColor(mtype)
		msg = fmt.Sprintf("%s%s[Peer %d]%s[%s] %s%s", BOLD, MAGENTA, id, msgcolor, mtype, RESET, msg)
	} else {
		msg = fmt.Sprintf("%s%s[Peer %d] %s%s", BOLD, MAGENTA, id, RESET, msg)
	}
	log.Printf("%s\n", msg)
	// fmt.Printf("%s\n", msg)
}

func ServerError(msg string, id int) {
	ServerLog(messages.ERROR, fmt.Sprintf("%sError: %s%s", RED, RESET, msg), id)
}

func GlobalLog(msg string) {
	log.Printf("%s%s[Node] %s%s\n", BOLD, BLUE, RESET, msg)
	// fmt.Printf("%s%s[Node] %s%s\n", BOLD, BLUE, RESET, msg)
}

func GlobalError(msg string) {
	GlobalLog(fmt.Sprintf("%sError: %s%s", RED, RESET, msg))
}
