package basedbhandler

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type DBHandler interface {
	Execute(ctx context.Context, statement string, args ...interface{}) error
	ExecuteBatch(ctx context.Context, statement string, args [][]interface{}) error
	Query(ctx context.Context, statement string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, statement string, args ...interface{}) (Row, error)
	Close()
}

type TransactionKey string

type TransactionalDBHandler interface {
	Execute(ctx context.Context, statement string, args ...interface{}) error
	ExecuteBatch(ctx context.Context, statement string, args [][]interface{}) error
	Query(ctx context.Context, statement string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, statement string, args ...interface{}) (Row, error)
	Transactioner
}

type Transactioner interface {
	Commit(ctx context.Context)
	Rollback(ctx context.Context)
	NewTx(ctx context.Context) (pgx.Tx, error)
}

type Rows interface {
	Scan(dest ...interface{}) error
	Next() bool
}

type Row interface {
	Scan(dest ...interface{}) error
}
