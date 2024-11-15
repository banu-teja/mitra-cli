package db

import (
	"database/sql"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *sql.DB
	*Queries
}

// NewStore creates a new store
func NewStore(connPool *sql.DB) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
