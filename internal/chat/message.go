package chat

type Message struct {
	// ClientID stores the id of the client that sent the message
	ClientID string
	// The data is the message transmited throught the websocket connection
	Data []byte
}
