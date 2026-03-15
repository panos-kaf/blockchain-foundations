package cli

import (
	"fmt"
	"marabu/internal/logs"
)

func Start() {
	fmt.Println("Marabu CLI started. Type 'help' for commands.")
	for {
		var cmd string
		fmt.Print("> ")
		fmt.Scanln(&cmd)

		switch cmd {
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  peers - List connected peers")
			fmt.Println("  objects - List stored objects")
			fmt.Println("  exit - Exit the CLI")
		case "peers":
			fmt.Print("dummy command")
			// logs.ClientLog("", fmt.Sprintf("Connected peers: %d", len(peer.GetConnectedPeers())), -1)
		case "objects":
			fmt.Print("dummy command")
			// logs.ClientLog("", fmt.Sprintf("Stored objects: %d", object.GetObjectCount()), -1)
		case "exit":
			logs.ClientLog("", "Exiting CLI...", -1)
			return
		default:
			logs.ClientError(fmt.Sprintf("Unknown command: %s", cmd), -1)
		}
	}
}
