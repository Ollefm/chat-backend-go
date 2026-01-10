package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/Ollefm/chat-backend-go/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepo struct {
	db *pgxpool.Pool
}

func NewChatRepo(db *pgxpool.Pool) *ChatRepo {
	return &ChatRepo{db}
}

func (r *ChatRepo) GetOrCreateChat(ctx context.Context, user1, user2 string) (*models.Chat, error) {
	if user1 == user2 {
		return nil, errors.New("cannot chat with yourself")
	}
	a, b := user1, user2
	if a > b {
		a, b = b, a
	}

	var chat models.Chat
	err := r.db.QueryRow(ctx,
		`SELECT c.id, c.user1_id, u1.username, c.user2_id, u2.username, c.created_at 
         FROM chats c
         JOIN users u1 ON c.user1_id = u1.id
         JOIN users u2 ON c.user2_id = u2.id
         WHERE c.user1_id=$1 AND c.user2_id=$2`,
		a, b).Scan(&chat.ID, &chat.User1ID, &chat.User1Username, &chat.User2ID, &chat.User2Username, &chat.Created)

	if err == nil {
		return &chat, nil
	}

	row := r.db.QueryRow(ctx,
		`INSERT INTO chats (user1_id, user2_id, created_at) VALUES ($1,$2,$3) 
         RETURNING id, user1_id, (SELECT username FROM users WHERE id=$1), user2_id, (SELECT username FROM users WHERE id=$2), created_at`,
		a, b, time.Now(),
	)
	if err := row.Scan(&chat.ID, &chat.User1ID, &chat.User1Username, &chat.User2ID, &chat.User2Username, &chat.Created); err != nil {
		return nil, err
	}

	return &chat, nil
}

func (r *ChatRepo) GetChatData(ctx context.Context, chatID string) (*models.Chat, error) {
	var chat models.Chat

	query := `SELECT c.id, c.user1_id, u1.username, c.user2_id, u2.username, c.created_at 
              FROM chats c
              JOIN users u1 ON c.user1_id = u1.id
              JOIN users u2 ON c.user2_id = u2.id
              WHERE c.id = $1`

	err := r.db.QueryRow(ctx, query, chatID).Scan(
		&chat.ID, &chat.User1ID, &chat.User1Username, &chat.User2ID, &chat.User2Username, &chat.Created,
	)

	if err != nil {
		return nil, err
	}

	return &chat, nil
}

func (r *ChatRepo) UserInChat(ctx context.Context, userID, chatID string) (bool, error) {
	var ok bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM chats WHERE id=$1 AND (user1_id=$2 OR user2_id=$2))`,
		chatID, userID).Scan(&ok)
	return ok, err
}

func (r *ChatRepo) GetParticipants(ctx context.Context, chatID string) ([]string, error) {
	rows, err := r.db.Query(ctx,
		`SELECT user1_id, user2_id FROM chats WHERE id=$1`,
		chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var u1, u2 string
		if err := rows.Scan(&u1, &u2); err != nil {
			return nil, err
		}
		ids = append(ids, u1, u2)
	}
	return ids, nil
}
