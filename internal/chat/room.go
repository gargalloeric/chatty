package chat

import "context"

type Room struct {
	// Map as a set of connected clients
	clients map[*Client]struct{}

	// Messages from the clients to be broadcasted
	Broadcast chan *Message

	// Clients joining the room
	Register chan *Client

	// Clients leaving the rooom
	Unregister chan *Client

	ctx context.Context

	// Close the room and disconnect all clients
	Shutdown context.CancelFunc
}

func NewRoom(ctx context.Context) *Room {
	ctx, cancel := context.WithCancel(ctx)
	return &Room{
		clients:    make(map[*Client]struct{}),
		Broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		ctx:        ctx,
		Shutdown:   cancel,
	}
}

func (h *Room) Run() {
	for {
		select {
		case <-h.ctx.Done():
			// Gracefully shutdown triggered, disconnect all clients
			for client := range h.clients {
				close(client.send)
				delete(h.clients, client)
			}
			return
		case client := <-h.Register:
			h.clients[client] = struct{}{}
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		// If we recieve a message, we have to send the message to every connected client
		case message := <-h.Broadcast:
			for client := range h.clients {
				if message.Sender != client.conn.RemoteAddr() {
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
}
