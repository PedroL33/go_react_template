package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgxDbConn struct {
	conn *pgx.Conn
}

// Ensure PgxDbConn implements DbConn
var _ DbConn = (*PgxDbConn)(nil)

// Begin implements DbConn
func (p *PgxDbConn) Begin(ctx context.Context) (DbConn, error) {
	tx, err := p.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &PgxTxWrapper{tx: tx}, nil
}

// Commit - Not applicable for PgxDbConn (but required by interface)
func (p *PgxDbConn) Commit(ctx context.Context) error {
	return nil // No-op for PgxDbConn
}

// Rollback - Not applicable for PgxDbConn
func (p *PgxDbConn) Rollback(ctx context.Context) error {
	return nil // No-op for PgxDbConn
}

// Exec implements DbConn
func (p *PgxDbConn) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return p.conn.Exec(ctx, sql, arguments...)
}

// Query implements DbConn
func (p *PgxDbConn) Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (DbRows, error) {
	return p.conn.Query(ctx, sql, optionsAndArgs...)
}

// QueryRow implements DbConn
func (p *PgxDbConn) QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) DbRow {
	return p.conn.QueryRow(ctx, sql, optionsAndArgs...)
}

// Conn returns itself (DbConn)
func (p *PgxDbConn) Conn() DbConn {
	return p
}
