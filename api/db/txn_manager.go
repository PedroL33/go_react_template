package db

import "context"

type TransactionManager interface {
	Begin(ctx context.Context) (DbConn, error)
}

type transactionManager struct {
	db DbConn
}

func NewTransactionManager(db DbConn) TransactionManager {
	return &transactionManager{
		db: db,
	}
}

func (t *transactionManager) Begin(ctx context.Context) (DbConn, error) {
	return t.db.Begin(ctx)
}

func GetDb(storeDb DbConn, tx DbConn) DbConn {
	if tx == nil {
		return storeDb
	}

	return tx
}
