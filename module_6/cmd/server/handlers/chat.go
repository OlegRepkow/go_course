package handlers

import (
	"log"
	"module_6/internal/services"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type ChatHandler struct {
	service *services.ChatService
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewChatHandler(service *services.ChatService) *ChatHandler {
	return &ChatHandler{
		service: service,
		clients: make(map[*websocket.Conn]bool),
	}
}

// History godoc
// @Summary Get chat history
// @Description Retrieve last 50 messages
// @Tags chat
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Message
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /channel/history [get]
func (h *ChatHandler) History(c *fiber.Ctx) error {
	username := c.Locals("username").(string)
	log.Printf("[INFO] History request from user=%s\n", username)

	messages, err := h.service.GetHistory(50)
	if err != nil {
		log.Printf("[ERROR] Failed to get history: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[INFO] Returned %d messages to user=%s\n", len(messages), username)
	return c.JSON(messages)
}

// Send godoc
// @Summary Send message
// @Description Send a message to the chat
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{text=string} true "Message text"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /channel/send [post]
func (h *ChatHandler) Send(c *fiber.Ctx) error {
	var req struct {
		Text string `json:"text"`
	}

	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ERROR] Send: invalid request body: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	username := c.Locals("username").(string)
	log.Printf("[INFO] Message from user=%s, text=%s\n", username, req.Text)

	if err := h.service.SaveMessage(username, req.Text); err != nil {
		log.Printf("[ERROR] Failed to save message from user=%s: %v\n", username, err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	h.broadcast(fiber.Map{"username": username, "text": req.Text})
	log.Printf("[INFO] Message broadcasted to %d clients\n", len(h.clients))
	return c.JSON(fiber.Map{"message": "sent"})
}

func (h *ChatHandler) Listen(c *websocket.Conn) {
	h.mu.Lock()
	h.clients[c] = true
	clientCount := len(h.clients)
	h.mu.Unlock()
	log.Printf("[INFO] WebSocket client connected, total clients: %d\n", clientCount)

	defer func() {
		h.mu.Lock()
		delete(h.clients, c)
		clientCount := len(h.clients)
		h.mu.Unlock()
		c.Close()
		log.Printf("[INFO] WebSocket client disconnected, total clients: %d\n", clientCount)
	}()

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			log.Printf("[ERROR] WebSocket read error: %v\n", err)
			break
		}
	}
}

func (h *ChatHandler) broadcast(msg interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for client := range h.clients {
		if err := client.WriteJSON(msg); err != nil {
			log.Printf("[ERROR] Failed to broadcast message: %v\n", err)
		}
	}
}
