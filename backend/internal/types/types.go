package types

import "time"

// Platform represents messaging platform type
type Platform string

const (
	PlatformWhatsApp  Platform = "whatsapp"
	PlatformInstagram Platform = "instagram"
	PlatformMessenger Platform = "messenger"
	PlatformWeb       Platform = "web"
)

// MessageDirection represents message direction
type MessageDirection string

const (
	DirectionInbound  MessageDirection = "inbound"
	DirectionOutbound MessageDirection = "outbound"
)

// MessageStatus represents message delivery status
type MessageStatus string

const (
	StatusPending   MessageStatus = "pending"
	StatusSent      MessageStatus = "sent"
	StatusDelivered MessageStatus = "delivered"
	StatusRead      MessageStatus = "read"
	StatusFailed    MessageStatus = "failed"
)

// Message represents a chat message
type Message struct {
	ID             string           `json:"id"`
	ConversationID string           `json:"conversation_id"`
	Platform       Platform         `json:"platform"`
	Direction      MessageDirection `json:"direction"`
	Content        string           `json:"content"`
	ContentType    string           `json:"content_type"` // text, image, video, etc
	Status         MessageStatus    `json:"status"`
	ExternalID     string           `json:"external_id"` // Meta message ID
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

// Conversation represents a chat conversation
type Conversation struct {
	ID              string    `json:"id"`
	ContactID       string    `json:"contact_id"`
	Platform        Platform  `json:"platform"`
	ExternalID      string    `json:"external_id"` // WhatsApp/IG thread ID
	LastMessageAt   time.Time `json:"last_message_at"`
	LastMessageText string    `json:"last_message_text"`
	UnreadCount     int       `json:"unread_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Joined data
	Contact  *Contact   `json:"contact,omitempty"`
	Messages []*Message `json:"messages,omitempty"`
}

// Contact represents a customer/contact
type Contact struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone,omitempty"`
	Email       string    `json:"email,omitempty"`
	WhatsAppID  string    `json:"whatsapp_id,omitempty"`
	InstagramID string    `json:"instagram_id,omitempty"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	Metadata    string    `json:"metadata,omitempty"` // JSON string for extra data
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SendMessageRequest represents outgoing message request
type SendMessageRequest struct {
	ConversationID string   `json:"conversation_id"`
	Platform       Platform `json:"platform"`
	RecipientID    string   `json:"recipient_id"` // Phone number or IG user ID
	Content        string   `json:"content"`
	ContentType    string   `json:"content_type"`
}

// WebhookPayload represents incoming Meta webhook
type WebhookPayload struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Time    int64  `json:"time"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
					ID        string `json:"id"`
					From      string `json:"from"`
					Timestamp string `json:"timestamp"`
					Type      string `json:"type"`
					Text      struct {
						Body string `json:"body"`
					} `json:"text"`
				} `json:"messages"`
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Recipient struct {
				ID string `json:"id"`
			} `json:"recipient"`
			Timestamp int64 `json:"timestamp"`
			Message   struct {
				Mid  string `json:"mid"`
				Text string `json:"text"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}
