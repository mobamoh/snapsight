package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
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
		UPDATE sessions
		SET token_hash = $2
		WHERE user_id = $1
    	RETURNING id;`, session.UserID, session.TokenHash)
	if err = row.Scan(&session.ID); errors.Is(err, sql.ErrNoRows) {
		row = r.DB.QueryRow(`
			INSERT INTO sessions (user_id, token_hash)
			VALUES ($1, $2)
			RETURNING id;`, session.UserID, session.TokenHash)
		err = row.Scan(&session.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}
func (r *SessionService) User(token string) (*User, error) {
	tokenHash := r.hash(token)
	var user User
	row := r.DB.QueryRow(`
		SELECT user_id FROM sessions WHERE token_hash = $1
	`, tokenHash)
	if err := row.Scan(&user.ID); err != nil {
		return nil, fmt.Errorf("session user: %w", err)
	}

	row = r.DB.QueryRow(`
		SELECT email, password_hash
		FROM users WHERE id = $1;`, user.ID)
	if err := row.Scan(&user.Email, &user.PasswordHash); err != nil {
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
