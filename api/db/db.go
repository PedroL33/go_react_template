package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Querier is the minimal interface required to run SQL. Both the connection
// pool and an in-flight transaction satisfy it, so stores can operate on
// either transparently.
type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...any) (DbRows, error)
	QueryRow(ctx context.Context, sql string, arguments ...any) DbRow
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

// Pool is a concurrency-safe wrapper around *pgxpool.Pool that satisfies the
// Querier interface using this package's row types.
type Pool struct {
	pool *pgxpool.Pool
}

// NewPool constructs a connection pool and verifies it with a ping.
func NewPool(ctx context.Context, databaseURL string) (*Pool, error) {
	p, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	if err := p.Ping(ctx); err != nil {
		p.Close()
		return nil, err
	}
	return &Pool{pool: p}, nil
}

func (p *Pool) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, sql, arguments...)
}

func (p *Pool) Query(ctx context.Context, sql string, arguments ...any) (DbRows, error) {
	return p.pool.Query(ctx, sql, arguments...)
}

func (p *Pool) QueryRow(ctx context.Context, sql string, arguments ...any) DbRow {
	return p.pool.QueryRow(ctx, sql, arguments...)
}

func (p *Pool) Close() {
	p.pool.Close()
}

// Acquire checks out a single connection from the pool. The caller MUST call
// Release on the returned Conn to return it to the pool; deferring the call
// immediately after acquiring is the safe pattern.
//
// Pinning a connection is useful for session-scoped work such as LISTEN/NOTIFY,
// session-level advisory locks, or temporary tables. For ordinary queries that
// don't require session state, use the pool directly.
func (p *Pool) Acquire(ctx context.Context) (*Conn, error) {
	c, err := p.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return &Conn{conn: c}, nil
}

// Conn is a single pinned connection from the pool. It satisfies Querier so
// it can be passed to stores like any other querier.
type Conn struct {
	conn *pgxpool.Conn
}

func (c *Conn) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return c.conn.Exec(ctx, sql, arguments...)
}

func (c *Conn) Query(ctx context.Context, sql string, arguments ...any) (DbRows, error) {
	return c.conn.Query(ctx, sql, arguments...)
}

func (c *Conn) QueryRow(ctx context.Context, sql string, arguments ...any) DbRow {
	return c.conn.QueryRow(ctx, sql, arguments...)
}

// Release returns the connection to the pool. Safe to call multiple times.
func (c *Conn) Release() {
	c.conn.Release()
}

// Compile-time checks that our wrappers satisfy Querier.
var (
	_ Querier = (*Pool)(nil)
	_ Querier = (*Conn)(nil)
)
