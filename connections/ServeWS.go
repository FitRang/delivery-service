package connections

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	email := strings.ToLower(strings.TrimSpace(r.Header.Get("x-user-email")))
	if email == "" {
		http.Error(w, "missing user-email", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	client := &Client{
		Email: email,
		Conn:  conn,
		Send:  make(chan []byte, 256),
	}

	hub.Register(client)

	go writePump(client)
	go readPump(client, hub)
}

func readPump(c *Client, hub *Hub) {
	defer func() {
		hub.Unregister(c)
		c.Conn.Close()
	}()

	for {
		if _, _, err := c.Conn.ReadMessage(); err != nil {
			break
		}
	}
}

func writePump(c *Client) {
	defer func() {
		log.Println("writePump closed for:", c.Email)
		c.Conn.Close()
	}()

	for msg := range c.Send {
		log.Printf("Writing WS message to %s: %s", c.Email, string(msg))

		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("WS write error:", err)
			return
		}
	}
}
