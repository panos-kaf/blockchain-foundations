//go:build no_bootstrap

package bootstrap

import (
	"marabu/internal/object"
	"marabu/internal/peer"
)

func StartNode(objectManager *object.ObjectManager) {
	go peer.StartServer(18018, objectManager)
}
