package db

import (
	"context"
	"database/sql"

	"github.com/openshift-online/rh-trex-core/db/transaction"
)

// By default do no roll back transaction.
// only perform rollback if explicitly set by g2.g2.MarkForRollback(ctx, err)
const defaultRollbackPolicy = false

// CoreSessionFactory interface for database session management in core library
type CoreSessionFactory interface {
	DirectDB() CoreDirectConnection
}

// CoreDirectConnection interface for direct database access in core library
type CoreDirectConnection interface {
	Begin() (*sql.Tx, error)
	QueryRow(query string, args ...interface{}) CoreRow
}

// CoreRow interface for database row operations in core library
type CoreRow interface {
	Scan(dest ...interface{}) error
}

// NewTransaction constructs a new Transaction object.
func NewTransaction(ctx context.Context, connection CoreSessionFactory) (*transaction.Transaction, error) {
	if connection == nil {
		// This happens in non-integration tests
		return nil, nil
	}

	dbx := connection.DirectDB()
	tx, err := dbx.Begin()
	if err != nil {
		return nil, err
	}

	// current transaction ID set by postgres.  these are *not* distinct across time
	// and do get reset after postgres performs "vacuuming" to reclaim used IDs.
	var txid int64
	row := tx.QueryRow("select txid_current()")
	if row != nil {
		err := row.Scan(&txid)
		if err != nil {
			return nil, err
		}
	}

	return transaction.Build(tx, txid, defaultRollbackPolicy), nil
}