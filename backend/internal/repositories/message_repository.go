package repositories

import (
	"context"
	"time"

	"github.com/temanbatin/omnichannel/internal/types"
)

type MessageRepository struct {
	db *DB
}

func NewMessageRepository(db *DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, msg *types.Message) error {
	query := `
		INSERT INTO messages (id, conversation_id, platform, direction, content, content_type, status, external_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		msg.ID, msg.ConversationID, msg.Platform, msg.Direction,
		msg.Content, msg.ContentType, msg.Status, msg.ExternalID,
		msg.CreatedAt, msg.UpdatedAt,
	)
	return err
}

func (r *MessageRepository) GetByID(ctx context.Context, id string) (*types.Message, error) {
	query := `
		SELECT id, conversation_id, platform, direction, content, content_type, status, external_id, created_at, updated_at
		FROM messages WHERE id = $1
	`
	msg := &types.Message{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&msg.ID, &msg.ConversationID, &msg.Platform, &msg.Direction,
		&msg.Content, &msg.ContentType, &msg.Status, &msg.ExternalID,
		&msg.CreatedAt, &msg.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (r *MessageRepository) ListByConversation(ctx context.Context, conversationID string, limit, offset int) ([]*types.Message, error) {
	query := `
		SELECT id, conversation_id, platform, direction, content, content_type, status, external_id, created_at, updated_at
		FROM messages 
		WHERE conversation_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*types.Message
	for rows.Next() {
		msg := &types.Message{}
		if err := rows.Scan(
			&msg.ID, &msg.ConversationID, &msg.Platform, &msg.Direction,
			&msg.Content, &msg.ContentType, &msg.Status, &msg.ExternalID,
			&msg.CreatedAt, &msg.UpdatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *MessageRepository) UpdateStatus(ctx context.Context, id string, status types.MessageStatus) error {
	query := `UPDATE messages SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Pool.Exec(ctx, query, status, time.Now(), id)
	return err
}
