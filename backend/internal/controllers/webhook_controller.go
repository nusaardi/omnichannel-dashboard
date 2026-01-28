package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/temanbatin/omnichannel/internal/config"
	"github.com/temanbatin/omnichannel/internal/services"
	"github.com/temanbatin/omnichannel/internal/types"
)

type WebhookController struct {
	messagingSvc *services.MessagingService
	config       *config.Config
}

func NewWebhookController(messagingSvc *services.MessagingService, cfg *config.Config) *WebhookController {
	return &WebhookController{
		messagingSvc: messagingSvc,
		config:       cfg,
	}
}

// VerifyWhatsApp handles WhatsApp webhook verification (GET)
func (c *WebhookController) VerifyWhatsApp(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	if mode == "subscribe" && token == c.config.MetaVerifyToken {
		log.Printf("WhatsApp webhook verified successfully")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
		return
	}

	log.Printf("WhatsApp webhook verification failed: mode=%s, token=%s", mode, token)
	http.Error(w, "Verification failed", http.StatusForbidden)
}

// HandleWhatsApp handles incoming WhatsApp webhooks (POST)
func (c *WebhookController) HandleWhatsApp(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read webhook body: %v", err)
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	// Verify signature if app secret is configured
	if c.config.MetaAppSecret != "" {
		signature := r.Header.Get("X-Hub-Signature-256")
		if !c.verifySignature(body, signature) {
			log.Printf("Invalid webhook signature")
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
	}

	var payload types.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Failed to parse webhook payload: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Process asynchronously to respond quickly
	go func() {
		if err := c.messagingSvc.ProcessIncomingWhatsApp(r.Context(), &payload); err != nil {
			log.Printf("Failed to process WhatsApp message: %v", err)
		}
	}()

	// Always respond 200 quickly to Meta
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("EVENT_RECEIVED"))
}

// VerifyInstagram handles Instagram webhook verification (GET)
func (c *WebhookController) VerifyInstagram(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	if mode == "subscribe" && token == c.config.MetaVerifyToken {
		log.Printf("Instagram webhook verified successfully")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
		return
	}

	log.Printf("Instagram webhook verification failed")
	http.Error(w, "Verification failed", http.StatusForbidden)
}

// HandleInstagram handles incoming Instagram webhooks (POST)
func (c *WebhookController) HandleInstagram(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read webhook body: %v", err)
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	// Verify signature if app secret is configured
	if c.config.MetaAppSecret != "" {
		signature := r.Header.Get("X-Hub-Signature-256")
		if !c.verifySignature(body, signature) {
			log.Printf("Invalid webhook signature")
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
	}

	var payload types.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Failed to parse webhook payload: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Process asynchronously
	go func() {
		if err := c.messagingSvc.ProcessIncomingInstagram(r.Context(), &payload); err != nil {
			log.Printf("Failed to process Instagram message: %v", err)
		}
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("EVENT_RECEIVED"))
}

// HandleWhatsAppInternal handles forwarded WhatsApp webhooks from n8n (no signature check)
func (c *WebhookController) HandleWhatsAppInternal(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read internal webhook body: %v", err)
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var payload types.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Failed to parse internal webhook payload: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Process asynchronously
	go func() {
		if err := c.messagingSvc.ProcessIncomingWhatsApp(r.Context(), &payload); err != nil {
			log.Printf("Failed to process WhatsApp message from n8n: %v", err)
		}
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// verifySignature verifies the Meta webhook signature
func (c *WebhookController) verifySignature(body []byte, signature string) bool {
	if signature == "" {
		return false
	}

	// Remove "sha256=" prefix
	if len(signature) > 7 && signature[:7] == "sha256=" {
		signature = signature[7:]
	}

	mac := hmac.New(sha256.New, []byte(c.config.MetaAppSecret))
	mac.Write(body)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
