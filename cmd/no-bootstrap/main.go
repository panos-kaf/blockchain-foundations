package main

import (
	"fmt"
	"marabu/internal/cli"
	"marabu/internal/logs"
	"marabu/internal/object"
	"marabu/internal/peer"
	"os"
	"path/filepath"
)

func main() {

	logFile := logs.InitLogs()
	defer logFile.Close()

	objectManager, err := object.NewObjectManager("./db")
	if err != nil {
		fmt.Printf("Error creating object manager: %v\n", err)
		return
	}

	PEERS_FILE := filepath.Join(".", "db", "peers.csv")
	peersFile, err := os.OpenFile(PEERS_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error creating peers file: %v\n", err)
		os.Exit(1)
	}
	defer peersFile.Close()

	go peer.StartServer(18018, objectManager)

	cli.Start()

}
