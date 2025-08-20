package main

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gargalloeric/chatty/internal/chat"
	"github.com/gargalloeric/chatty/internal/component/view"
	"github.com/gorilla/websocket"
)

func waitForMessage(sub <-chan chat.Message) tea.Cmd {
	return func() tea.Msg {
		message := <-sub
		switch message.Type {
		case chat.TextType:
			var payload chat.Text
			json.Unmarshal(message.Payload, &payload)
			return view.TextMsg{Text: payload.Content, Sender: "Anonymous"}
		case chat.MetadataType:
			var payload chat.Metadata
			json.Unmarshal(message.Payload, &payload)
			return view.MetadataMsg{Room: payload.Room, UserCount: payload.UserCount}
		default:
			return nil
		}
	}
}

func writeToConn(conn *websocket.Conn, message string) tea.Cmd {
	return func() tea.Msg {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			return errorMsg(err)
		}

		return view.TextMsg{Text: message, Sender: "You"}
	}
}

func listenFromConn(conn *websocket.Conn, sub chan<- chat.Message) tea.Cmd {
	return func() tea.Msg {
		var data chat.Message
		for {

			if err := conn.ReadJSON(&data); err != nil {
				return errorMsg(err)
			}
			sub <- data
		}
	}
}
