package chat

type Hub struct {
	// Map as a set of connected clients
	clients map[*Client]struct{}

	// Messages from the clients to be broadcasted
	Broadcast chan []byte

	// Clients joining the hub
	Register chan *Client

	// Clients leaving the hub
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
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
