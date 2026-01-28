package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/temanbatin/omnichannel/internal/services"
	"github.com/temanbatin/omnichannel/internal/types"
)

type MessageController struct {
	messagingSvc *services.MessagingService
}

func NewMessageController(messagingSvc *services.MessagingService) *MessageController {
	return &MessageController{messagingSvc: messagingSvc}
}

// List returns all messages (paginated)
func (c *MessageController) List(w http.ResponseWriter, r *http.Request) {
	// For now, return empty - messages are accessed via conversations
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"messages": []interface{}{},
		"message":  "Use /api/conversations/{id} to get messages",
	})
}

// Get returns a single message
func (c *MessageController) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	respondJSON(w, http.StatusOK, map[string]string{"id": id})
}

// Send sends a new message
func (c *MessageController) Send(w http.ResponseWriter, r *http.Request) {
	var req types.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Content == "" {
		respondError(w, http.StatusBadRequest, "Content is required")
		return
	}

	if req.Platform == "" {
		respondError(w, http.StatusBadRequest, "Platform is required")
		return
	}

	msg, err := c.messagingSvc.SendMessage(r.Context(), &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, msg)
}

// ListConversations returns all conversations
func (c *MessageController) ListConversations(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0

	conversations, err := c.messagingSvc.ListConversations(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"conversations": conversations,
		"total":         len(conversations),
	})
}

// GetConversation returns a conversation with messages
func (c *MessageController) GetConversation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	conversation, err := c.messagingSvc.GetConversation(r.Context(), id, 50)
	if err != nil {
		respondError(w, http.StatusNotFound, "Conversation not found")
		return
	}

	respondJSON(w, http.StatusOK, conversation)
}

// ListContacts returns all contacts
func (c *MessageController) ListContacts(w http.ResponseWriter, r *http.Request) {
	limit := 100
	offset := 0

	contacts, err := c.messagingSvc.ListContacts(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"contacts": contacts,
		"total":    len(contacts),
	})
}

// CreateContact creates a new contact
func (c *MessageController) CreateContact(w http.ResponseWriter, r *http.Request) {
	var contact types.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if contact.Name == "" {
		respondError(w, http.StatusBadRequest, "Name is required")
		return
	}

	if err := c.messagingSvc.CreateContact(r.Context(), &contact); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, contact)
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
