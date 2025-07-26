package chat

import "net"

type Message struct {
	// The sender represents the client address
	Sender net.Addr
	// The data is the message transmited throught the websocket connection
	Data []byte
}
