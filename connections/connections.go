package connections

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Email string
	Conn   *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	clients map[string]map[*Client]bool
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[c.Email]; !ok {
		h.clients[c.Email] = make(map[*Client]bool)
	}
	h.clients[c.Email][c] = true
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if conns, ok := h.clients[c.Email]; ok {
		delete(conns, c)
		if len(conns) == 0 {
			delete(h.clients, c.Email)
		}
	}
}

func (h *Hub) SendToUser(email string, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conns, ok := h.clients[email]
	if !ok {
		log.Println("No active WS connection for:", email)
		return
	}

	log.Printf("Delivering message to %s (%d connections)", email, len(conns))

	for c := range conns {
		select {
		case c.Send <- msg:
		default:
			log.Println("Send buffer full for:", email)
		}
	}
}
