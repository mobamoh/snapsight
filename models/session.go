package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/mobamoh/snapsight/rand"
)

const (
	MinBytesPerToken = 32
)

type Session struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
}

type SessionService struct {
	DB            *sql.DB
	BytesPerToken int
}

func (r *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := r.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("session create: %w", err)
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: r.hash(token),
	}

	row := r.DB.QueryRow(`
			INSERT INTO sessions (user_id, token_hash)
				VALUES ($1, $2) 
			ON CONFLICT (user_id) DO
				UPDATE SET token_hash = $2
			RETURNING id;`, session.UserID, session.TokenHash)
	if err = row.Scan(&session.ID); err != nil {
		return nil, fmt.Errorf("session create: %w", err)
	}

	return &session, nil
}
func (r *SessionService) User(token string) (*User, error) {
	tokenHash := r.hash(token)
	var user User
	row := r.DB.QueryRow(`
		SELECT u.id, u.email, u.password_hash
		FROM sessions s JOIN users u on u.id = s.user_id 
		WHERE s.token_hash = $1;`, tokenHash)

	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		return nil, fmt.Errorf("session user: %w", err)
	}
	return &user, nil
}
func (r *SessionService) Delete(token string) error {
	tokenHash := r.hash(token)
	_, err := r.DB.Exec(`
 		DELETE FROM sessions WHERE token_hash = $1;
 	`, tokenHash)
	if err != nil {
		return fmt.Errorf("session delete: %w", err)
	}
	return nil
}

func (r *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
