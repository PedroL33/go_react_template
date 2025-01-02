package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DbConn interface {
	Begin(ctx context.Context) (DbConn, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (DbRows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) DbRow
	Conn() DbConn
	Close(ctx context.Context) error
}

func NewDbConn(ctx context.Context, databaseURL string) (DbConn, error) {
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	return &PgxDbConn{conn: conn}, nil
}

type DbRow interface {
	Scan(dest ...any) error
}

type DbRows interface {
	Close()
	Err() error
	Scan(dest ...any) error
	Values() ([]any, error)
	Next() bool
}
