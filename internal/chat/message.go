package chat

type Message struct {
	// From stores the id of the client that sent the message
	From string `json:"from"`

	// Text represents the message text payload encoded as UTF-8 characters
	Text string `json:"text"`
}
