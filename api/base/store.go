package base

import (
	"example/dashboard/api/db"
)

type Store struct {
	DB db.DbConn
}

func NewBaseStore(db db.DbConn) *Store {
	return &Store{
		DB: db,
	}
}
