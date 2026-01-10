package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

/* Client represents a websocket connection, having two different purposes:

1. To receive messages from the client.
2. To send messages to the client.

For example User A and User B are both connected to ChatID 123.

User A sends a messsage "hello" ->
Client A Reads the messsage ->
The hub (separate logic) recieves it and processes it ->
Client B writes to User B the "hello" message.

That is broadcasting messages to all clients with the same ChatID.
*/

type Client struct {
	conn   *websocket.Conn
	userID string
	chatID string
	send   chan []byte
}

func NewClient(conn *websocket.Conn, userID, chatID string) *Client {
	return &Client{
		conn:   conn,
		userID: userID,
		chatID: chatID,
		send:   make(chan []byte, 256), // Buffered channel
	}
}

func (c *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.Unregister(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(512)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		hub.Broadcast(&Message{
			ChatID: c.chatID,
			UserID: c.userID,
			Data:   message,
		})
	}
}

func (c *Client) WriteMessage() {
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}
