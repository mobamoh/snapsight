package models

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/mobamoh/snapsight/errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

var (
	// A common pattern is to add the package as a prefix to the error for context.
	ErrEmailTaken = errors.New("models: email address is already in use")
)

type User struct {
	ID           int //TODO: change to uid
	Email        string
	PasswordHash string
	// TODO: add more fields
}

type UserService struct {
	DB *sql.DB
}

func (service *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	user := User{
		Email:        email,
		PasswordHash: string(hashedPwd),
	}

	row := service.DB.QueryRow(`INSERT INTO users (email, password_hash)
	VALUES ($1, $2) RETURNING id`, user.Email, user.PasswordHash)
	if err = row.Scan(&user.ID); err != nil {
		// See if we can use this error as a PgError
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			// This is a PgError, so see if it matches a unique violation.
			if pgError.Code == pgerrcode.UniqueViolation {
				// If this is true, it has to be an email violation since this is the
				// only way to trigger this type of violation with our SQL.
				return nil, ErrEmailTaken
			}
		}
		return nil, fmt.Errorf("user create: %w", err)
	}
	return &user, nil
}

func (service *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := User{
		Email: email,
	}

	row := service.DB.QueryRow(`SELECT id, password_hash FROM users WHERE email = $1`, email)
	if err := row.Scan(&user.ID, &user.PasswordHash); err != nil {
		return nil, fmt.Errorf("user authenticate: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("user authenticate: %w", err)
	}
	return &user, nil
}

func (service *UserService) UpdatePassword(userID int, password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	passwordHash := string(hashedBytes)
	_, err = service.DB.Exec(`
	  UPDATE users
		SET password_hash = $2
		WHERE id = $1;`, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}
