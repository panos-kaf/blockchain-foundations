package main

import (
	"fmt"
	"marabu/internal/cli"
	"marabu/internal/logs"
	"marabu/internal/object"
	"marabu/internal/peer"
	"net"
	"os"
	"path/filepath"
	"strconv"
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

	for _, p := range peer.BOOTSTRAP_PEERS {
		go func(p string) {
			host, portStr, _ := net.SplitHostPort(p)
			port, _ := strconv.Atoi(portStr)
			err := peer.StartClient(host, port, objectManager, func() {
				// Handle client disconnect if needed
			})
			if err != nil {
				fmt.Printf("Error connecting to peer %s: %v\n", p, err)
			}
		}(p)
	}

	cli.Start()

}
