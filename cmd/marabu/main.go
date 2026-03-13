package main

import (
	"fmt"
	"marabu/internal/peer"
	"net"
	"strconv"
)

func main() {

	go peer.StartServer(18018)

	for _, p := range peer.BOOTSTRAP_PEERS {
		go func(p string) {
			host, portStr, _ := net.SplitHostPort(p)
			port, _ := strconv.Atoi(portStr)
			err := peer.StartClient(host, port, 0, func() {
				// Handle client disconnect if needed
			})
			fmt.Println("Connected to peer: " + p)
			if err != nil {
				fmt.Printf("Error connecting to peer %s: %v\n", p, err)
			}
		}(p)
	}

	select {}

}
