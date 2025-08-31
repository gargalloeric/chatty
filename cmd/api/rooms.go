package main

import (
	"context"
	"encoding/json"
	"maps"
	"net/http"
	"slices"

	"github.com/gargalloeric/chatty/internal/chat"
)

func (app *application) listRoomHandler(w http.ResponseWriter, r *http.Request) {
	defer app.mx.RUnlock()
	w.Header().Set("Content-Type", "application/json")

	app.mx.RLock()
	rooms := slices.Collect(maps.Values(app.rooms))
	if err := json.NewEncoder(w).Encode(rooms); err != nil {
		http.Error(w, "encoding failed", http.StatusInternalServerError)
		return
	}
}

func (app *application) createRoomHandler(w http.ResponseWriter, r *http.Request) {
	defer app.mx.Unlock()

	var input struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "payload decoding failed", http.StatusBadRequest)
		return
	}

	app.mx.Lock()
	room := chat.NewRoom(context.Background(), app.logger, input.Name)
	if _, ok := app.rooms[room.ID]; ok {
		http.Error(w, "invalid room id", http.StatusInternalServerError)
		return
	}
	app.background(func() {
		room.Run()
	})

	app.rooms[room.ID] = room

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(room); err != nil {
		http.Error(w, "encoding failed", http.StatusInternalServerError)
	}
}

func (app *application) connectToRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("roomID")
	if roomID == "" {
		http.Error(w, "invalid roomID", http.StatusBadRequest)
		return
	}

	app.mx.RLock()
	room, ok := app.rooms[roomID]
	if !ok {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}
	app.mx.RUnlock()

	conn, err := app.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := chat.NewClient(room, conn, app.logger)

	room.Register <- client

	app.background(func() {
		client.Read()
	})

	app.background(func() {
		client.Write()
	})
}
