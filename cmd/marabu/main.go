package main

import (
	"fmt"
	"marabu/internal/cli"
	"marabu/internal/object"
	"marabu/internal/peer"
	"marabu/internal/logs"
	"net"
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
