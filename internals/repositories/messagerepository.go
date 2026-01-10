package repositories

import (
	"context"
	"time"

	"github.com/Ollefm/chat-backend-go/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepo struct {
	db *pgxpool.Pool
}

func NewMessageRepo(db *pgxpool.Pool) *MessageRepo {
	return &MessageRepo{db}
}

func (r *MessageRepo) SaveMessage(ctx context.Context, m *models.Message) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO messages (chat_id, sender_id, content, created_at) VALUES ($1, $2, $3, $4)`,
		m.ChatID, m.SenderID, m.Content, time.Now(),
	)
	return err
}

func (r *MessageRepo) GetMessages(ctx context.Context, chatID string) ([]models.Message, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT id, chat_id, sender_id, content, created_at 
		 FROM messages 
		 WHERE chat_id=$1 
		 ORDER BY created_at DESC LIMIT 20`,
		chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.ChatID, &m.SenderID, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}

	return msgs, nil
}
