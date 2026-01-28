package repositories

import (
	"context"
	"time"

	"github.com/temanbatin/omnichannel/internal/types"
)

type ConversationRepository struct {
	db *DB
}

func NewConversationRepository(db *DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) Create(ctx context.Context, conv *types.Conversation) error {
	query := `
		INSERT INTO conversations (id, contact_id, platform, external_id, last_message_at, last_message_text, unread_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		conv.ID, conv.ContactID, conv.Platform, conv.ExternalID,
		conv.LastMessageAt, conv.LastMessageText, conv.UnreadCount,
		conv.CreatedAt, conv.UpdatedAt,
	)
	return err
}

func (r *ConversationRepository) GetByID(ctx context.Context, id string) (*types.Conversation, error) {
	query := `
		SELECT c.id, c.contact_id, c.platform, c.external_id, c.last_message_at, c.last_message_text, c.unread_count, c.created_at, c.updated_at,
		       ct.id, ct.name, ct.phone, ct.email, ct.whatsapp_id, ct.instagram_id, ct.avatar_url
		FROM conversations c
		LEFT JOIN contacts ct ON c.contact_id = ct.id
		WHERE c.id = $1
	`
	conv := &types.Conversation{Contact: &types.Contact{}}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&conv.ID, &conv.ContactID, &conv.Platform, &conv.ExternalID,
		&conv.LastMessageAt, &conv.LastMessageText, &conv.UnreadCount,
		&conv.CreatedAt, &conv.UpdatedAt,
		&conv.Contact.ID, &conv.Contact.Name, &conv.Contact.Phone,
		&conv.Contact.Email, &conv.Contact.WhatsAppID, &conv.Contact.InstagramID,
		&conv.Contact.AvatarURL,
	)
	if err != nil {
		return nil, err
	}
	return conv, nil
}

func (r *ConversationRepository) GetByContactAndPlatform(ctx context.Context, contactID string, platform types.Platform) (*types.Conversation, error) {
	query := `
		SELECT id, contact_id, platform, external_id, last_message_at, last_message_text, unread_count, created_at, updated_at
		FROM conversations
		WHERE contact_id = $1 AND platform = $2
	`
	conv := &types.Conversation{}
	err := r.db.Pool.QueryRow(ctx, query, contactID, platform).Scan(
		&conv.ID, &conv.ContactID, &conv.Platform, &conv.ExternalID,
		&conv.LastMessageAt, &conv.LastMessageText, &conv.UnreadCount,
		&conv.CreatedAt, &conv.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return conv, nil
}

func (r *ConversationRepository) List(ctx context.Context, limit, offset int) ([]*types.Conversation, error) {
	query := `
		SELECT c.id, c.contact_id, c.platform, c.external_id, c.last_message_at, c.last_message_text, c.unread_count, c.created_at, c.updated_at,
		       ct.id, ct.name, ct.phone, ct.avatar_url
		FROM conversations c
		LEFT JOIN contacts ct ON c.contact_id = ct.id
		ORDER BY c.last_message_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*types.Conversation
	for rows.Next() {
		conv := &types.Conversation{Contact: &types.Contact{}}
		if err := rows.Scan(
			&conv.ID, &conv.ContactID, &conv.Platform, &conv.ExternalID,
			&conv.LastMessageAt, &conv.LastMessageText, &conv.UnreadCount,
			&conv.CreatedAt, &conv.UpdatedAt,
			&conv.Contact.ID, &conv.Contact.Name, &conv.Contact.Phone,
			&conv.Contact.AvatarURL,
		); err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}
	return conversations, nil
}

func (r *ConversationRepository) UpdateLastMessage(ctx context.Context, id, messageText string) error {
	query := `
		UPDATE conversations 
		SET last_message_at = $1, last_message_text = $2, unread_count = unread_count + 1, updated_at = $1
		WHERE id = $3
	`
	_, err := r.db.Pool.Exec(ctx, query, time.Now(), messageText, id)
	return err
}

func (r *ConversationRepository) MarkAsRead(ctx context.Context, id string) error {
	query := `UPDATE conversations SET unread_count = 0, updated_at = $1 WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, time.Now(), id)
	return err
}
