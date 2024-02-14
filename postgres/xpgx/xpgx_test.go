package xpgx

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/crossworth/pkgs/postgres/test"
	"github.com/crossworth/pkgs/xerror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestWithinTransaction(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	testDB := test.NewPostgresTest(t,
		test.WithDeleteDatabaseFunction(test.ForceDeleteDatabaseFunction),
		test.WithBaseAddress("postgres://postgres:root@localhost:1432"),
	)
	pool, err := pgxpool.New(ctx, testDB)
	require.NoError(t, err)
	// check if we have a connection on the context, should be false
	require.False(t, IsWithinTransaction(ctx))
	err = WithinTransaction(ctx, pool, func(ctx context.Context) error {
		// check if we are inside a transaction, should be true
		require.True(t, IsWithinTransaction(ctx))
		// try to start another transaction
		err = WithinTransaction(ctx, pool, func(ctx context.Context) error {
			// check if we are inside a transaction, should be true
			require.True(t, IsWithinTransaction(ctx))
			return nil
		})
		require.NoError(t, err)
		// create a table
		_, err = ConnectionFromContext(ctx).Exec(ctx, `CREATE TABLE a (b VARCHAR);`)
		require.NoError(t, err)
		// insert a few records
		_, err = ConnectionFromContext(ctx).Exec(ctx, `INSERT INTO a VALUES ('a'), ('b'), ('c'), ('d');`)
		require.NoError(t, err)
		// check if the records are created
		rows, _ := ConnectionFromContext(ctx).Query(ctx, `SELECT COUNT(*) FROM a`)
		count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
		require.NoError(t, err)
		require.Equal(t, 4, count)
		// return an error to trigger a rollback
		return fmt.Errorf("some error")
	})
	require.Error(t, err)
	rows, _ := pool.Query(ctx, `SELECT COUNT(*) FROM a`)
	_, err = pgx.CollectOneRow(rows, pgx.RowTo[int])
	require.Error(t, err)
	require.ErrorContains(t, err, "SQLSTATE 42P01")
}

func TestHandleError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	testDB := test.NewPostgresTest(t,
		test.WithDeleteDatabaseFunction(test.ForceDeleteDatabaseFunction),
		test.WithBaseAddress("postgres://postgres:root@localhost:1432"),
	)
	pool, err := pgxpool.New(ctx, testDB)
	require.NoError(t, err)
	// create a table
	_, err = pool.Exec(ctx, `CREATE TABLE a (b VARCHAR);`)
	require.NoError(t, err)
	// add a unique constraint
	_, err = pool.Exec(ctx, `CREATE UNIQUE INDEX test_idx ON a (b);`)
	require.NoError(t, err)
	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		err := pool.QueryRow(ctx, `SELECT * FROM a WHERE b = 'test' LIMIT 1;`).Scan()
		require.Error(t, err)
		err = HandleError("a", err)
		var xError xerror.Error
		require.True(t, errors.As(err, &xError))
		require.Equal(t, "not_found: entity=a", err.Error())
	})
	t.Run("constraint error", func(t *testing.T) {
		t.Parallel()
		// add a record
		_, err := pool.Exec(ctx, `INSERT INTO a VALUES ('a');`)
		require.NoError(t, err)
		// try to insert a second time
		_, err = pool.Exec(ctx, `INSERT INTO a VALUES ('a');`)
		require.Error(t, err)
		err = HandleError("a", err)
		var xError xerror.Error
		require.True(t, errors.As(err, &xError))
		require.Equal(t, "bad_request: constraint=test_idx, entity=a", err.Error())
	})
}
