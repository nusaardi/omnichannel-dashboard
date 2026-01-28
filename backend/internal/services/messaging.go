package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/temanbatin/omnichannel/internal/config"
	"github.com/temanbatin/omnichannel/internal/repositories"
	"github.com/temanbatin/omnichannel/internal/types"
	"github.com/temanbatin/omnichannel/pkg/meta"
)

// MessagingService handles unified messaging across platforms
type MessagingService struct {
	messageRepo      *repositories.MessageRepository
	contactRepo      *repositories.ContactRepository
	conversationRepo *repositories.ConversationRepository
	config           *config.Config

	whatsappClient  *meta.WhatsAppClient
	instagramClient *meta.InstagramClient
}

// NewMessagingService creates a new messaging service
func NewMessagingService(
	messageRepo *repositories.MessageRepository,
	contactRepo *repositories.ContactRepository,
	conversationRepo *repositories.ConversationRepository,
	cfg *config.Config,
) *MessagingService {
	svc := &MessagingService{
		messageRepo:      messageRepo,
		contactRepo:      contactRepo,
		conversationRepo: conversationRepo,
		config:           cfg,
	}

	// Initialize platform clients
	if cfg.MetaAccessToken != "" && cfg.WhatsAppPhoneID != "" {
		svc.whatsappClient = meta.NewWhatsAppClient(
			cfg.MetaAccessToken,
			cfg.WhatsAppPhoneID,
			cfg.WhatsAppBusinessID,
		)
	}

	if cfg.MetaAccessToken != "" && cfg.InstagramAccountID != "" {
		svc.instagramClient = meta.NewInstagramClient(
			cfg.MetaAccessToken,
			cfg.InstagramAccountID,
		)
	}

	return svc
}

// SendMessage sends a message to a recipient
func (s *MessagingService) SendMessage(ctx context.Context, req *types.SendMessageRequest) (*types.Message, error) {
	now := time.Now()
	msg := &types.Message{
		ID:             uuid.New().String(),
		ConversationID: req.ConversationID,
		Platform:       req.Platform,
		Direction:      types.DirectionOutbound,
		Content:        req.Content,
		ContentType:    req.ContentType,
		Status:         types.StatusPending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// Send via platform API
	var externalID string
	var err error

	switch req.Platform {
	case types.PlatformWhatsApp:
		if s.whatsappClient == nil {
			return nil, fmt.Errorf("WhatsApp client not configured")
		}
		resp, sendErr := s.whatsappClient.SendText(req.RecipientID, req.Content)
		if sendErr != nil {
			msg.Status = types.StatusFailed
			s.messageRepo.Create(ctx, msg)
			return nil, fmt.Errorf("failed to send WhatsApp message: %w", sendErr)
		}
		if len(resp.Messages) > 0 {
			externalID = resp.Messages[0].ID
		}

	case types.PlatformInstagram:
		if s.instagramClient == nil {
			return nil, fmt.Errorf("Instagram client not configured")
		}
		resp, sendErr := s.instagramClient.SendText(req.RecipientID, req.Content)
		if sendErr != nil {
			msg.Status = types.StatusFailed
			s.messageRepo.Create(ctx, msg)
			return nil, fmt.Errorf("failed to send Instagram message: %w", sendErr)
		}
		externalID = resp.MessageID

	default:
		return nil, fmt.Errorf("unsupported platform: %s", req.Platform)
	}

	msg.ExternalID = externalID
	msg.Status = types.StatusSent

	// Save message
	if err = s.messageRepo.Create(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	// Update conversation
	s.conversationRepo.UpdateLastMessage(ctx, req.ConversationID, req.Content)

	return msg, nil
}

// ProcessIncomingWhatsApp processes incoming WhatsApp webhook
func (s *MessagingService) ProcessIncomingWhatsApp(ctx context.Context, payload *types.WebhookPayload) error {
	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			if change.Field != "messages" {
				continue
			}

			for _, waMsg := range change.Value.Messages {
				// Get or create contact
				contact, err := s.getOrCreateWhatsAppContact(ctx, waMsg.From, change.Value.Contacts)
				if err != nil {
					return fmt.Errorf("failed to get/create contact: %w", err)
				}

				// Get or create conversation
				conversation, err := s.getOrCreateConversation(ctx, contact.ID, types.PlatformWhatsApp)
				if err != nil {
					return fmt.Errorf("failed to get/create conversation: %w", err)
				}

				// Create message
				now := time.Now()
				msg := &types.Message{
					ID:             uuid.New().String(),
					ConversationID: conversation.ID,
					Platform:       types.PlatformWhatsApp,
					Direction:      types.DirectionInbound,
					Content:        waMsg.Text.Body,
					ContentType:    waMsg.Type,
					Status:         types.StatusDelivered,
					ExternalID:     waMsg.ID,
					CreatedAt:      now,
					UpdatedAt:      now,
				}

				if err := s.messageRepo.Create(ctx, msg); err != nil {
					return fmt.Errorf("failed to save message: %w", err)
				}

				// Update conversation
				s.conversationRepo.UpdateLastMessage(ctx, conversation.ID, waMsg.Text.Body)
			}
		}
	}

	return nil
}

// ProcessIncomingInstagram processes incoming Instagram webhook
func (s *MessagingService) ProcessIncomingInstagram(ctx context.Context, payload *types.WebhookPayload) error {
	for _, entry := range payload.Entry {
		for _, messaging := range entry.Messaging {
			senderID := messaging.Sender.ID

			// Get or create contact
			contact, err := s.getOrCreateInstagramContact(ctx, senderID)
			if err != nil {
				return fmt.Errorf("failed to get/create contact: %w", err)
			}

			// Get or create conversation
			conversation, err := s.getOrCreateConversation(ctx, contact.ID, types.PlatformInstagram)
			if err != nil {
				return fmt.Errorf("failed to get/create conversation: %w", err)
			}

			// Create message
			now := time.Now()
			msg := &types.Message{
				ID:             uuid.New().String(),
				ConversationID: conversation.ID,
				Platform:       types.PlatformInstagram,
				Direction:      types.DirectionInbound,
				Content:        messaging.Message.Text,
				ContentType:    "text",
				Status:         types.StatusDelivered,
				ExternalID:     messaging.Message.Mid,
				CreatedAt:      now,
				UpdatedAt:      now,
			}

			if err := s.messageRepo.Create(ctx, msg); err != nil {
				return fmt.Errorf("failed to save message: %w", err)
			}

			// Update conversation
			s.conversationRepo.UpdateLastMessage(ctx, conversation.ID, messaging.Message.Text)
		}
	}

	return nil
}

// ListConversations returns all conversations
func (s *MessagingService) ListConversations(ctx context.Context, limit, offset int) ([]*types.Conversation, error) {
	return s.conversationRepo.List(ctx, limit, offset)
}

// GetConversation returns a conversation with messages
func (s *MessagingService) GetConversation(ctx context.Context, id string, messageLimit int) (*types.Conversation, error) {
	conv, err := s.conversationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	messages, err := s.messageRepo.ListByConversation(ctx, id, messageLimit, 0)
	if err != nil {
		return nil, err
	}

	conv.Messages = messages
	return conv, nil
}

// ListContacts returns all contacts
func (s *MessagingService) ListContacts(ctx context.Context, limit, offset int) ([]*types.Contact, error) {
	return s.contactRepo.List(ctx, limit, offset)
}

// CreateContact creates a new contact
func (s *MessagingService) CreateContact(ctx context.Context, contact *types.Contact) error {
	now := time.Now()
	contact.ID = uuid.New().String()
	contact.CreatedAt = now
	contact.UpdatedAt = now
	return s.contactRepo.Create(ctx, contact)
}

// Helper: get or create WhatsApp contact
func (s *MessagingService) getOrCreateWhatsAppContact(ctx context.Context, waID string, waContacts []struct {
	Profile struct {
		Name string `json:"name"`
	} `json:"profile"`
	WaID string `json:"wa_id"`
}) (*types.Contact, error) {
	contact, err := s.contactRepo.GetByWhatsAppID(ctx, waID)
	if err == nil {
		return contact, nil
	}

	// Create new contact
	name := waID
	for _, c := range waContacts {
		if c.WaID == waID && c.Profile.Name != "" {
			name = c.Profile.Name
			break
		}
	}

	now := time.Now()
	contact = &types.Contact{
		ID:         uuid.New().String(),
		Name:       name,
		Phone:      waID,
		WhatsAppID: waID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.contactRepo.Create(ctx, contact); err != nil {
		return nil, err
	}

	return contact, nil
}

// Helper: get or create Instagram contact
func (s *MessagingService) getOrCreateInstagramContact(ctx context.Context, igID string) (*types.Contact, error) {
	contact, err := s.contactRepo.GetByInstagramID(ctx, igID)
	if err == nil {
		return contact, nil
	}

	// Try to get profile from Instagram
	name := igID
	if s.instagramClient != nil {
		profile, profErr := s.instagramClient.GetUserProfile(igID)
		if profErr == nil {
			if n, ok := profile["name"].(string); ok && n != "" {
				name = n
			} else if u, ok := profile["username"].(string); ok && u != "" {
				name = u
			}
		}
	}

	now := time.Now()
	contact = &types.Contact{
		ID:          uuid.New().String(),
		Name:        name,
		InstagramID: igID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.contactRepo.Create(ctx, contact); err != nil {
		return nil, err
	}

	return contact, nil
}

// Helper: get or create conversation
func (s *MessagingService) getOrCreateConversation(ctx context.Context, contactID string, platform types.Platform) (*types.Conversation, error) {
	conv, err := s.conversationRepo.GetByContactAndPlatform(ctx, contactID, platform)
	if err == nil {
		return conv, nil
	}

	now := time.Now()
	conv = &types.Conversation{
		ID:            uuid.New().String(),
		ContactID:     contactID,
		Platform:      platform,
		LastMessageAt: now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.conversationRepo.Create(ctx, conv); err != nil {
		return nil, err
	}

	return conv, nil
}
