package main

import (
	"bytes"
	"net/http"
	"time"

	ws "github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	hub  *Hub
	conn *ws.Conn
	send chan []byte
}

func (c *Client) read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(appData string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		message = bytes.TrimSpace(bytes.ReplaceAll(message, newline, space))
		c.hub.broadcast <- message
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel, so we inform the peer that we are closing the connection
				c.conn.WriteMessage(ws.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(ws.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			for i := 0; i < len(c.send); i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type Hub struct {
	// Map as a set of connected clients
	clients map[*Client]struct{}

	// Messages from the clients to be broadcasted
	broadcast chan []byte

	// Clients joining the hub
	register chan *Client

	// Clients leaving the hub
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		// If we recieve a message, we have to send the message to every connected client
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				// If we cannot send the message, we assumed that the client is dead or stuck
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}

		}
	}
}

func (app *application) chatHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := app.upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	client := &Client{hub: app.hub, conn: conn, send: make(chan []byte)}
	app.hub.register <- client

	go client.read()
	go client.write()
}
