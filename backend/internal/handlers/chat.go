package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/therizza/freakdream/backend/internal/middleware"
	"github.com/therizza/freakdream/backend/internal/models"
	"github.com/therizza/freakdream/backend/internal/repository"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type client struct {
	conn   *websocket.Conn
	userID int64
	send   chan []byte
	ctx    context.Context
}

type ChatHub struct {
	mu      sync.RWMutex
	clients map[int64]*client
	chat    *repository.ChatRepository
}

func NewChatHub(chat *repository.ChatRepository) *ChatHub {
	return &ChatHub{
		clients: make(map[int64]*client),
		chat:    chat,
	}
}

func (h *ChatHub) ServeWS(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	c := &client{conn: conn, userID: userID, send: make(chan []byte, 256), ctx: r.Context()}

	h.mu.Lock()
	h.clients[userID] = c
	h.mu.Unlock()

	go h.writePump(c)
	h.readPump(c)
}

func (h *ChatHub) readPump(c *client) {
	defer func() {
		h.mu.Lock()
		delete(h.clients, c.userID)
		h.mu.Unlock()
		c.conn.Close()
	}()

	c.conn.SetReadLimit(4096)
	_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg models.WSMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "message":
			if msg.ReceiverID == 0 || msg.Texto == "" {
				continue
			}
			_, _ = h.chat.Save(c.ctx, c.userID, msg.ReceiverID, msg.Texto)

			out, _ := json.Marshal(models.WSMessage{
				Type:     "message",
				SenderID: c.userID,
				Texto:    msg.Texto,
				Data:     time.Now().Format(time.RFC3339),
			})

			h.mu.RLock()
			if receiver, ok := h.clients[msg.ReceiverID]; ok {
				select {
				case receiver.send <- out:
				default:
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *ChatHub) writePump(c *client) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// History returns recent messages between two users.
func (h *ChatHub) History(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	otherID, err := parseIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	msgs, err := h.chat.GetConversation(r.Context(), userID, otherID, 50)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to load messages")
		return
	}

	respondJSON(w, http.StatusOK, msgs)
}
