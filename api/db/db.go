package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -source=db.go -destination=../users/mocks/mock_querier.go -package=mocks

// Querier is the minimal interface required to run SQL. Both the connection
// pool and an in-flight transaction satisfy it, so stores can operate on
// either transparently.
type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row
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

// NewPool constructs a connection pool and verifies it with a ping.
func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	p, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	if err := p.Ping(ctx); err != nil {
		p.Close()
		return nil, err
	}
	return p, nil
}
