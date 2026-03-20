//go:build standard

package bootstrap

import (
	"fmt"
	"marabu/internal/messages"
	"marabu/internal/object"
	"marabu/internal/peer"
	"net"
	"strconv"
)

// start server and initiate client connections to bootstrap peers
func StartNode(objectManager *object.ObjectManager) {
	go peer.StartServer(18018, objectManager)

	for _, p := range peer.BOOTSTRAP_PEERS {
		go func(p messages.Peer) {
			host, portStr, _ := net.SplitHostPort(string(p))
			port, _ := strconv.Atoi(portStr)
			err := peer.StartClient(host, port, objectManager, func() {
				// Handle client disconnect if needed
			})
			if err != nil {
				fmt.Printf("Error connecting to peer %s: %v\n", p, err)
			}
		}(p)
	}
}
