package main

import (
	"net/http"

	"github.com/gargalloeric/chatty/internal/chat"
)

func (app *application) chatHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := app.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := chat.NewClient(app.room, conn)
	app.room.Register <- client

	app.background(func() {
		client.Read()
	})

	app.background(func() {
		client.Write()
	})
}
