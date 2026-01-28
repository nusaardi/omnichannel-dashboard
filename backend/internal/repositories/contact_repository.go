package repositories

import (
	"context"
	"time"

	"github.com/temanbatin/omnichannel/internal/types"
)

type ContactRepository struct {
	db *DB
}

func NewContactRepository(db *DB) *ContactRepository {
	return &ContactRepository{db: db}
}

func (r *ContactRepository) Create(ctx context.Context, contact *types.Contact) error {
	query := `
		INSERT INTO contacts (id, name, phone, email, whatsapp_id, instagram_id, avatar_url, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		contact.ID, contact.Name, contact.Phone, contact.Email,
		contact.WhatsAppID, contact.InstagramID, contact.AvatarURL,
		contact.Metadata, contact.CreatedAt, contact.UpdatedAt,
	)
	return err
}

func (r *ContactRepository) GetByID(ctx context.Context, id string) (*types.Contact, error) {
	query := `
		SELECT id, name, phone, email, whatsapp_id, instagram_id, avatar_url, metadata, created_at, updated_at
		FROM contacts WHERE id = $1
	`
	contact := &types.Contact{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&contact.ID, &contact.Name, &contact.Phone, &contact.Email,
		&contact.WhatsAppID, &contact.InstagramID, &contact.AvatarURL,
		&contact.Metadata, &contact.CreatedAt, &contact.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (r *ContactRepository) GetByWhatsAppID(ctx context.Context, waID string) (*types.Contact, error) {
	query := `
		SELECT id, name, phone, email, whatsapp_id, instagram_id, avatar_url, metadata, created_at, updated_at
		FROM contacts WHERE whatsapp_id = $1
	`
	contact := &types.Contact{}
	err := r.db.Pool.QueryRow(ctx, query, waID).Scan(
		&contact.ID, &contact.Name, &contact.Phone, &contact.Email,
		&contact.WhatsAppID, &contact.InstagramID, &contact.AvatarURL,
		&contact.Metadata, &contact.CreatedAt, &contact.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (r *ContactRepository) GetByInstagramID(ctx context.Context, igID string) (*types.Contact, error) {
	query := `
		SELECT id, name, phone, email, whatsapp_id, instagram_id, avatar_url, metadata, created_at, updated_at
		FROM contacts WHERE instagram_id = $1
	`
	contact := &types.Contact{}
	err := r.db.Pool.QueryRow(ctx, query, igID).Scan(
		&contact.ID, &contact.Name, &contact.Phone, &contact.Email,
		&contact.WhatsAppID, &contact.InstagramID, &contact.AvatarURL,
		&contact.Metadata, &contact.CreatedAt, &contact.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (r *ContactRepository) List(ctx context.Context, limit, offset int) ([]*types.Contact, error) {
	query := `
		SELECT id, name, phone, email, whatsapp_id, instagram_id, avatar_url, metadata, created_at, updated_at
		FROM contacts
		ORDER BY updated_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []*types.Contact
	for rows.Next() {
		contact := &types.Contact{}
		if err := rows.Scan(
			&contact.ID, &contact.Name, &contact.Phone, &contact.Email,
			&contact.WhatsAppID, &contact.InstagramID, &contact.AvatarURL,
			&contact.Metadata, &contact.CreatedAt, &contact.UpdatedAt,
		); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

func (r *ContactRepository) Update(ctx context.Context, contact *types.Contact) error {
	query := `
		UPDATE contacts 
		SET name = $1, phone = $2, email = $3, whatsapp_id = $4, instagram_id = $5, 
		    avatar_url = $6, metadata = $7, updated_at = $8
		WHERE id = $9
	`
	_, err := r.db.Pool.Exec(ctx, query,
		contact.Name, contact.Phone, contact.Email, contact.WhatsAppID,
		contact.InstagramID, contact.AvatarURL, contact.Metadata,
		time.Now(), contact.ID,
	)
	return err
}
