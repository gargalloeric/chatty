package chat

import (
	"bytes"
	"context"
	"time"

	"github.com/gargalloeric/chatty/internal/identity"
	"github.com/gorilla/websocket"
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
	id   string
	hub  *Room
	conn *websocket.Conn
	send chan *Message
	ctx  context.Context
}

func NewClient(hub *Room, conn *websocket.Conn) *Client {
	id, err := identity.GenerateRandomID(16)
	if err != nil {
		panic("unable to generate client id")
	}

	return &Client{
		id:   id,
		hub:  hub,
		conn: conn,
		send: make(chan *Message),
		ctx:  hub.ctx,
	}
}

func (c *Client) Read() {
	defer c.conn.Close()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(appData string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				c.hub.Unregister <- c
			}
			break
		}
		message = bytes.TrimSpace(bytes.ReplaceAll(message, newline, space))
		c.hub.Broadcast <- &Message{Sender: c.conn.RemoteAddr(), Data: message}
	}
}

func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case <-c.ctx.Done():
			// Cancellation signal received so we exit the event loop
			return
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel, so we inform the peer that we are closing the connection
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message.Data)

			// Add queued messages to the current websocket message
			for i := 0; i < len(c.send); i++ {
				w.Write(newline)

				message := <-c.send
				w.Write(message.Data)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
