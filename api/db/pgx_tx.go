package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgxTxWrapper struct {
	tx pgx.Tx
}

// Ensure PgxTxWrapper implements DbConn
var _ DbConn = (*PgxTxWrapper)(nil)

// Begin is not applicable for PgxTxWrapper
func (p *PgxTxWrapper) Begin(ctx context.Context) (DbConn, error) {
	return nil, errors.New("cannot begin a transaction within a transaction")
}

func (p *PgxTxWrapper) Close(ctx context.Context) error {
	return errors.New("transaction cannot be closed")
}

// Commit implements DbConn
func (p *PgxTxWrapper) Commit(ctx context.Context) error {
	return p.tx.Commit(ctx)
}

// Rollback implements DbConn
func (p *PgxTxWrapper) Rollback(ctx context.Context) error {
	return p.tx.Rollback(ctx)
}

// Exec implements DbConn
func (p *PgxTxWrapper) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return p.tx.Exec(ctx, sql, arguments...)
}

// Query implements DbConn
func (p *PgxTxWrapper) Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (DbRows, error) {
	return p.tx.Query(ctx, sql, optionsAndArgs...)
}

// QueryRow implements DbConn
func (p *PgxTxWrapper) QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) DbRow {
	return p.tx.QueryRow(ctx, sql, optionsAndArgs...)
}

// Conn returns itself (DbConn)
func (p *PgxTxWrapper) Conn() DbConn {
	return p
}
