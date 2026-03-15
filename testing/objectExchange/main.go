package main

import (
	"bufio"
	"fmt"
	"marabu/internal/messages"
	"net"
	"os"
)

func send(conn net.Conn, msg string) {
	fmt.Fprintf(conn, "%s\n", msg)
}

func receive(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	resp, _ := reader.ReadString('\n')
	return resp
}

func exchangeObject(objectID messages.HashID, objectMsg string, conn net.Conn, resp string) {
	// 1. Send ihaveobject
	ihaveMsg, _ := messages.MakeIHaveObjectMessage(objectID)
	send(conn, ihaveMsg)
	fmt.Println("Sent ihaveobject")

	// 2. Expect getobject
	resp = receive(conn)
	fmt.Println("Received:", resp)
	// Parse and check for getobject

	// 3. Send object
	send(conn, objectMsg)
	fmt.Println("Sent object")

	// 4. Expect ihaveobject gossip (optional, if you have multiple peers)
	resp = receive(conn)
	fmt.Println("Received:", resp)

	// 5. Send getobject for known object
	getObjMsg, _ := messages.MakeGetObjectMessage(objectID)
	send(conn, getObjMsg)
	fmt.Println("Sent getobject")

	// 6. Expect object response
	resp = receive(conn)
	fmt.Println("Received:", resp)

}

func main() {
	serverAddr := "localhost:18018" // Change to your server address
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Failed to connect:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Example object ID and object (replace with real values)
	objectID := messages.HashID("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	// fmt.Printf("Object ID = %s\n", objectID)

	dummyTransaction := messages.Transaction{
		Type: messages.TRANSACTION,
		Inputs: []messages.TxInput{
			{
				Outpoint: messages.Outpoint{
					Txid:  "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
					Index: 0,
				},
				Sig: "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
			},
		},
		Outputs: []messages.TxOutput{
			{
				Pubkey: "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
				Value:  100,
			},
		},
	}

	validID, _ := messages.HashObject(dummyTransaction)

	objectMsg, _ := messages.MakeTXObjectMessage(dummyTransaction)

	invalidObject := messages.ObjectSchema{
		Type:     messages.OBJECT,
		ObjectID: objectID, //invalid id
		Object:   dummyTransaction,
	}

	// fmt.Printf("Generated object message: %s\n", objectMsg)

	// 0. Greet the server
	helloMsg, _ := messages.MakeHelloMessage()
	send(conn, helloMsg)
	fmt.Println("Sent hello")

	resp := receive(conn)
	fmt.Println("Received:", resp)
	// Parse and check for hello response

	fmt.Println("Starting 1st object exchange...")
	exchangeObject(validID, objectMsg, conn, resp)

	fmt.Println("Starting 2nd object exchange with invalid object...")
	// Try sending an invalid object
	invalidObjectMsg, _ := messages.CanonicalizeMessage(invalidObject)
	exchangeObject(objectID, invalidObjectMsg, conn, resp)
}
