//go:build no_bootstrap

package main

import (
	"marabu/internal/cli"
	"marabu/internal/object"
	"marabu/internal/peer"
)

// start server but dont initiate any client connections to bootstrap peers
func startNode(objectManager *object.ObjectManager) {

	go peer.StartServer(18018, objectManager)

	cli.Start()

}
