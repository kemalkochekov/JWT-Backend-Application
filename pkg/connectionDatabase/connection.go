package connectionDatabase

import (
	"Fiber_JWT_Authentication_backend_server/configs"
	"context"
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBops interface {
	// database queries
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Close() error
}

type Database struct {
	db *sqlx.DB
}

func NewDatabase(inDB *sqlx.DB) *Database {
	return &Database{
		db: inDB,
	}
}
func (s *Database) GetPool() *sqlx.DB {
	return s.db
}
func (s *Database) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.db.GetContext(ctx, dest, query, args...)
}
func (s *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}
func (s *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}
func (s *Database) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.db.SelectContext(ctx, dest, query, args...)
}
func (s *Database) Close() error {
	if err := goose.Down(s.db.DB, "./internal/repository/migrations"); err != nil {
		fmt.Printf("goose migration down failed: %v", err)
	}
	return s.db.Close()
}
func GenerateDsn(cfgs configs.DatabaseConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfgs.Host, cfgs.Port, cfgs.User, cfgs.Password, cfgs.DBName)
}

// cfgs database configuration from env file
func NewDB(ctx context.Context, cfgs configs.DatabaseConfig) (*Database, error) {
	db, err := sqlx.Connect("postgres", GenerateDsn(cfgs))
	if err != nil {
		return nil, fmt.Errorf("could not create connection pool: %v", err)
	}
	if err := goose.Up(db.DB, "./internal/repository/migrations"); err != nil {
		return nil, fmt.Errorf("goose migration up failed: %v", err)
	}
	return &Database{db: db}, nil
}
