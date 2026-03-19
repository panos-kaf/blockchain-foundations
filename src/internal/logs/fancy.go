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
	WHITE   = "\033[37m"

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

func Log(mtype messages.MessageType, msg string, id int, color string) {
	if mtype != "" {
		msgcolor := MessageTypeColor(mtype)
		msg = fmt.Sprintf("%s%s[%d]%s[%s] %s%s", BOLD, color, id, msgcolor, mtype, RESET, msg)
	} else {
		msg = fmt.Sprintf("%s%s[%d] %s%s", BOLD, color, id, RESET, msg)
	}
	log.Printf("%s\n", msg)
	// fmt.Printf("%s\n", msg)
}

func Error(mtype messages.MessageType, msg string, id int, color string) {
	Log(mtype, fmt.Sprintf("%sError: %s%s", RED, msg, RESET), id, color)
}

func Message(mtype messages.MessageType, sends bool, id int, color string, addr string) {
	var direction string
	if sends {
		direction = CYAN + "[-->]" + RESET
	} else {
		direction = YELLOW + "[<--]" + RESET
	}
	msgcolor := MessageTypeColor(mtype)

	// example format: [Peer 1][HELLO][-->][127.0.0.1:12345]
	log.Printf("%s%s[%d]%s[%s]%s%s[%s]%s", BOLD, color, id, msgcolor, mtype, direction, WHITE, addr, RESET)
}

func MessageError(mtype messages.MessageType, msg string, sends bool, id int, color string, addr string) {
	var direction string
	if sends {
		direction = CYAN + "[-->]" + RESET
	} else {
		direction = YELLOW + "[<--]" + RESET
	}
	msgcolor := MessageTypeColor(mtype)

	log.Printf("%s%s[%d]%s[%s]%s%s[%s] %sError: %s%s", BOLD, color, id, msgcolor, mtype, direction, WHITE, addr, RED, msg, RESET)
}

// Peer-specific log functions

func ClientLog(mtype messages.MessageType, msg string, id int) {
	Log(mtype, msg, id, BLUE)
}

func ClientError(mtype messages.MessageType, msg string, id int) {
	Error(mtype, msg, id, BLUE)
}

func ServerLog(mtype messages.MessageType, msg string, id int) {
	Log(mtype, msg, id, MAGENTA)
}

func ClientMessage(mtype messages.MessageType, sends bool, id int, addr string) {
	Message(mtype, sends, id, BLUE, addr)
}

func ClientMessageError(mtype messages.MessageType, msg string, sends bool, id int, addr string) {
	MessageError(mtype, msg, sends, id, BLUE, addr)
}

func ServerMessage(mtype messages.MessageType, sends bool, id int, addr string) {
	Message(mtype, sends, id, MAGENTA, addr)
}

func ServerMessageError(mtype messages.MessageType, msg string, sends bool, id int, addr string) {
	MessageError(mtype, msg, sends, id, MAGENTA, addr)
}

func ServerError(mtype messages.MessageType, msg string, id int) {
	Error(mtype, msg, id, MAGENTA)
}

// Global log functions for non-peer-specific messages

func GlobalLog(msg string) {
	log.Printf("%s%s[*] %s%s\n", BOLD, BLUE, RESET, msg)
	// fmt.Printf("%s%s[Node] %s%s\n", BOLD, BLUE, RESET, msg)
}

func GlobalError(msg string) {
	GlobalLog(fmt.Sprintf("%sError: %s%s", RED, RESET, msg))
}
