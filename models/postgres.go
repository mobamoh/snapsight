package models

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	"io/fs"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "admin",
		Password: "nimda",
		Database: "snapsight",
		SSLMode:  "disable",
	}
}

func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("postgres open: %w", err)
	}
	return db, nil
}

func Migrate(db *sql.DB, dir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("postgres migrate: %v", err)
	}
	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("postgres migrate: %v", err)
	}
	return nil
}

func MigrateFS(db *sql.DB, fs fs.FS, dir string) error {
	goose.SetBaseFS(fs)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}
