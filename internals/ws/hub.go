package ws

type Message struct {
	ChatID string
	UserID string
	Data   []byte
}

type Hub struct {
	clients       map[*Client]string
	clientsByUser map[string]*Client
	broadcast     chan *Message
	register      chan *Client
	unregister    chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:       make(map[*Client]string),
		clientsByUser: make(map[string]*Client),
		broadcast:     make(chan *Message),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:

			if existingClient, exists := h.clientsByUser[client.userID]; exists {
				delete(h.clients, existingClient)
				if existingClient.send != nil {
					close(existingClient.send)
				}
			}
			h.clients[client] = client.userID
			h.clientsByUser[client.userID] = client

		case client := <-h.unregister:

			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.clientsByUser, client.userID)
				close(client.send)
			}

		case message := <-h.broadcast:
			var targets []*Client
			for client := range h.clients {
				if client.chatID == message.ChatID {
					targets = append(targets, client)
				}
			}

			for _, client := range targets {
				select {
				case client.send <- message.Data:
				default:
					h.Unregister(client)
				}
			}
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Broadcast(message *Message) {
	h.broadcast <- message
}
