// Package xpgx provides a few interfaces
// and helper functions to work with pgx.
package xpgx

import (
	"context"
	"errors"
	"strings"

	"github.com/crossworth/pkgs/xerror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Queryable interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

type Executable interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

type Txable interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Acquirable interface {
	Acquire(ctx context.Context) (*pgx.Conn, error)
}

type Connection interface {
	Queryable
	Executable
	Txable
}

type connectionKey struct{}

// ConnectionFromContext returns a Connection from the given context.
// It will return nil if no Connection is found in the context.
func ConnectionFromContext(ctx context.Context) Connection {
	queryableExecutable, ok := ctx.Value(connectionKey{}).(Connection)
	if !ok {
		return nil
	}
	return queryableExecutable
}

// SetConnectionOnContext sets the Connection on the context.
func SetConnectionOnContext(ctx context.Context, connection Connection) context.Context {
	return context.WithValue(ctx, connectionKey{}, connection)
}

type transactionKey struct{}

// IsWithinTransaction returns true if we are inside a transaction.
func IsWithinTransaction(ctx context.Context) bool {
	inside, ok := ctx.Value(transactionKey{}).(bool)
	return ok && inside
}

// WithinTransaction executes the given function inside a transaction.
// It's important to note that we do not support nested transactions (save points).
func WithinTransaction(ctx context.Context, conn Connection, inTransaction func(ctx context.Context) error) error {
	// already in a transaction
	if IsWithinTransaction(ctx) {
		return inTransaction(ctx)
	}
	return pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		ctx = SetConnectionOnContext(ctx, tx)
		ctx = context.WithValue(ctx, transactionKey{}, true)
		return inTransaction(ctx)
	})
}

// HandleError handle the Postgres errors converting for our repository errors.
func HandleError(typ string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return xerror.MakeNotFoundError(xerror.ErrParam("entity", typ))
	}
	if strings.Contains(err.Error(), "SQLSTATE 23505") {
		var (
			idx    = strings.Index(err.Error(), `constraint "`)
			params = []xerror.ErrorParam{
				xerror.ErrParam("entity", typ),
			}
		)
		if idx != -1 {
			sub := err.Error()[idx+12:]
			if lastIdx := strings.Index(sub, `"`); lastIdx != -1 {
				sub = sub[:lastIdx]
			}
			params = append(params, xerror.ErrParam("constraint", sub))
		}
		return xerror.MakeBadRequestError(params...)
	}
	return err
}

// EnsureAffected check if the pgconn.CommandTag has the number of records affected,
// otherwise we create a not found error.
func EnsureAffected(typ string, res pgconn.CommandTag, n int64) error {
	if res.RowsAffected() == n {
		return nil
	}
	return xerror.MakeNotFoundError(xerror.ErrParam("entity", typ))
}
