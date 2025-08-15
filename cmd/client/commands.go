package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gargalloeric/chatty/internal/chat"
	"github.com/gargalloeric/chatty/internal/component/view"
	"github.com/gorilla/websocket"
)

func waitForMessage(sub <-chan chat.Message) tea.Cmd {
	return func() tea.Msg {
		message := <-sub
		return view.TextMsg{Text: message.Text, Sender: "Anonymous"}
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
