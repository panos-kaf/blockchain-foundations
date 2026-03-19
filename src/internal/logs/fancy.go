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
	case messages.MSG_HELLO:
		return GREEN
	case messages.MSG_ERROR:
		return RED
	case messages.MSG_GETPEERS, messages.MSG_PEERS:
		return CYAN
	case messages.MSG_GETOBJECT, messages.MSG_IHAVEOBJECT, messages.MSG_OBJECT:
		return YELLOW
	case messages.MSG_GETMEMPOOL, messages.MSG_MEMPOOL:
		return BLUE
	case messages.MSG_GETCHAINTIP, messages.MSG_CHAINTIP:
		return MAGENTA
	default:
		return RESET
	}
}

func Log(mtype messages.MessageType, msg string, id int, color string) {
	if mtype != messages.MSG_NONE {
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

func Message(mtype messages.MessageType, errorCode messages.ErrorCode, sends bool, id int, color string, addr string) {
	var direction string
	if sends {
		direction = CYAN + "[-->]" + RESET
	} else {
		direction = YELLOW + "[<--]" + RESET
	}
	msgcolor := MessageTypeColor(mtype)

	label := string(mtype)
	if mtype == messages.MSG_ERROR && errorCode != messages.E_NONE {
		label = string(errorCode)
	}

	// example format: [Peer 1][MSG_HELLO][-->][127.0.0.1:12345]
	log.Printf("%s%s[%d]%s[%s]%s%s[%s]%s", BOLD, color, id, msgcolor, label, direction, WHITE, addr, RESET)
}

func MessageError(mtype messages.MessageType, errorCode messages.ErrorCode, msg string, sends bool, id int, color string, addr string) {
	var direction string
	if sends {
		direction = CYAN + "[-->]" + RESET
	} else {
		direction = YELLOW + "[<--]" + RESET
	}
	msgcolor := MessageTypeColor(mtype)

	label := string(mtype)
	if mtype == messages.MSG_ERROR && errorCode != messages.E_NONE {
		label = string(errorCode)
	}
	log.Printf("%s%s[%d]%s[%s]%s%s[%s] %sError: %s%s", BOLD, color, id, msgcolor, label, direction, WHITE, addr, RED, msg, RESET)
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

func ClientMessage(mtype messages.MessageType, errorCode messages.ErrorCode, sends bool, id int, addr string) {
	Message(mtype, errorCode, sends, id, BLUE, addr)
}

func ClientMessageError(mtype messages.MessageType, errorCode messages.ErrorCode, msg string, sends bool, id int, addr string) {
	MessageError(mtype, errorCode, msg, sends, id, BLUE, addr)
}

func ServerMessage(mtype messages.MessageType, errorCode messages.ErrorCode, sends bool, id int, addr string) {
	Message(mtype, errorCode, sends, id, MAGENTA, addr)
}

func ServerMessageError(mtype messages.MessageType, errorCode messages.ErrorCode, msg string, sends bool, id int, addr string) {
	MessageError(mtype, errorCode, msg, sends, id, MAGENTA, addr)
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
