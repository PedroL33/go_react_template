package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -source=txn_manager.go -destination=../users/mocks/mock_txn_manager.go -package=mocks

// TransactionManager runs a function within a database transaction, handling
// begin, commit, and rollback automatically. Rollback is invoked if fn returns
// an error or panics; commit runs only on success.
type TransactionManager interface {
	WithTx(ctx context.Context, fn func(tx Querier) error) error
}

type transactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) TransactionManager {
	return &transactionManager{pool: pool}
}

func (t *transactionManager) WithTx(ctx context.Context, fn func(tx Querier) error) (err error) {
	pgxTx, err := t.pool.Begin(ctx)
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

	if err = fn(pgxTx); err != nil {
		return err
	}

	if err = pgxTx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
