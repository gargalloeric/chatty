package chat

import "encoding/json"

type MessageType int

const (
	TextType MessageType = iota
	MetadataType
)

type Message struct {
	// Type specifies whether the message is a text message or other type of message
	Type MessageType `json:"type"`

	// Payload contains the content as a json RawMessage type
	Payload json.RawMessage `json:"payload"`
}

type Text struct {
	From    string `json:"from"`
	Content string `json:"content"`
}

type Metadata struct {
	Room      string `json:"room"`
	UserCount int    `json:"user_count"`
}
