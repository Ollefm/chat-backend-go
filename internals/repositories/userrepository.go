package repositories

import (
	"context"
	"time"

	"github.com/Ollefm/chat-backend-go/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) CreateUser(ctx context.Context, username, passwordHash string) (string, error) {

	dbContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var userID string
	err := r.db.QueryRow(dbContext,
		`INSERT INTO users (username, password_hash) 
		 VALUES ($1, $2) 
		 RETURNING id`,
		username, passwordHash).Scan(&userID)

	if err != nil {
		return "", err
	}
	return userID, nil
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (string, string, error) {

	dbContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var userID, passwordHash string
	err := r.db.QueryRow(dbContext,
		`SELECT id, password_hash 
		 FROM users 
		 WHERE username = $1`,
		username).Scan(&userID, &passwordHash)

	if err != nil {
		return "", "", err
	}
	return userID, passwordHash, nil
}

func (r *UserRepo) GetUsernameByUid(ctx context.Context, uid string) (string, error) {

	dbContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var username string
	err := r.db.QueryRow(dbContext,
		`SELECT username 
     FROM users 
     WHERE id = $1`,
		uid).Scan(&username)

	if err != nil {
		return "", err
	}
	return username, nil
}

func (r *UserRepo) SearchUsers(ctx context.Context, query string, excludeID string) ([]models.UserResult, error) {
	dbContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(dbContext,
		`SELECT id, username
         FROM users
         WHERE username ILIKE $1
           AND id != $2
         LIMIT 5`,
		"%"+query+"%", excludeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]models.UserResult, 0)

	for rows.Next() {
		var u models.UserResult
		if err := rows.Scan(&u.ID, &u.Username); err != nil {
			return nil, err
		}
		results = append(results, u)
	}

	return results, nil
}
