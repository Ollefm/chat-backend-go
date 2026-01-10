package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepo struct {
	db *pgxpool.Pool
}

func NewSessionRepo(db *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{db}
}

// CreateSession inserts a new session into the database
func (r *SessionRepo) CreateSession(ctx context.Context, sessionHash, userID string, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx,
		`INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)`,
		sessionHash, userID, expiresAt)
	return err
}

// GetSession retrieves a session and returns the associated user ID and expiry
func (r *SessionRepo) GetSession(ctx context.Context, sessionHash string) (string, time.Time, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var userID string
	var expiresAt time.Time

	err := r.db.QueryRow(ctx,
		`SELECT user_id, expires_at FROM sessions WHERE id = $1`,
		sessionHash).Scan(&userID, &expiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", time.Time{}, errors.New("session not found")
		}
		return "", time.Time{}, err
	}

	return userID, expiresAt, nil
}

// DeleteSession deletes a single session by hash
func (r *SessionRepo) DeleteSession(ctx context.Context, sessionHash string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, sessionHash)
	return err
}

// DeleteAllUserSessions deletes all sessions for a given user
func (r *SessionRepo) DeleteAllUserSessions(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID)
	return err
}
