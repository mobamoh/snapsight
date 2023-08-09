package models

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	ID           uint //TODO: change to uid
	Email        string
	PasswordHash string
	// TODO: add more fields
}

type UserService struct {
	DB *sql.DB
}

func (u *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("user create: %w", err)
	}
	user := User{
		Email:        email,
		PasswordHash: string(hashedPwd),
	}

	row := u.DB.QueryRow(`INSERT INTO users (email, password_hash)
	VALUES ($1, $2) RETURNING id`, user.Email, user.PasswordHash)
	if err = row.Scan(&user.ID); err != nil {
		return nil, fmt.Errorf("user create: %w", err)
	}
	return &user, nil
}

func (u *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := User{
		Email: email,
	}

	row := u.DB.QueryRow(`SELECT id, password_hash FROM users WHERE email = $1`, email)
	if err := row.Scan(&user.ID, &user.PasswordHash); err != nil {
		return nil, fmt.Errorf("user authenticate: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("user authenticate: %w", err)
	}
	return &user, nil
}
