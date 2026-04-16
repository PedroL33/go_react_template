package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// TransactionManager runs a function within a database transaction, handling
// begin, commit, and rollback automatically. Rollback is invoked if fn returns
// an error or panics; commit runs only on success.
type TransactionManager interface {
	WithTx(ctx context.Context, fn func(tx Querier) error) error
}

type transactionManager struct {
	pool *Pool
}

func NewTransactionManager(pool *Pool) TransactionManager {
	return &transactionManager{pool: pool}
}

func (t *transactionManager) WithTx(ctx context.Context, fn func(tx Querier) error) (err error) {
	pgxTx, err := t.pool.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = pgxTx.Rollback(ctx)
			panic(p)
		}
		if err != nil {
			if rbErr := pgxTx.Rollback(ctx); rbErr != nil && rbErr != pgx.ErrTxClosed {
				err = fmt.Errorf("rollback after error (%w): %v", err, rbErr)
			}
		}
	}()

	if err = fn(&txQuerier{tx: pgxTx}); err != nil {
		return err
	}

	if err = pgxTx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

// txQuerier adapts a pgx.Tx to the Querier interface.
type txQuerier struct {
	tx pgx.Tx
}

func (q *txQuerier) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return q.tx.Exec(ctx, sql, arguments...)
}

func (q *txQuerier) Query(ctx context.Context, sql string, arguments ...any) (DbRows, error) {
	return q.tx.Query(ctx, sql, arguments...)
}

func (q *txQuerier) QueryRow(ctx context.Context, sql string, arguments ...any) DbRow {
	return q.tx.QueryRow(ctx, sql, arguments...)
}

var _ Querier = (*txQuerier)(nil)
