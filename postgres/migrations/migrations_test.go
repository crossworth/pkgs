package migrations

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/crossworth/pkgs/postgres/test"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func Test_migrationPrefix(t *testing.T) {
	t.Parallel()
	require.Equal(t, "001", migrationPrefix("001.sql"))
	require.Equal(t, "001", migrationPrefix("001.tx.sql"))
	require.Equal(t, "001", migrationPrefix("001.no-tx.sql"))
	require.Equal(t, "001", migrationPrefix("001.no-tx.other.sql"))
}

func Test_sortMigrations(t *testing.T) {
	t.Parallel()
	t.Run("simple case", func(t *testing.T) {
		t.Parallel()
		sorted, err := sortMigrations([]string{"002.sql", "001.sql", "003.sql"})
		require.NoError(t, err)
		require.Equal(t, []string{"001.sql", "002.sql", "003.sql"}, sorted)
	})
	t.Run("tx/no-txt", func(t *testing.T) {
		t.Parallel()
		sorted, err := sortMigrations([]string{"002.no-tx.sql", "001.sql", "003.tx.sql"})
		require.NoError(t, err)
		require.Equal(t, []string{"001.sql", "002.no-tx.sql", "003.tx.sql"}, sorted)
	})
}

func Test_buildMigrationPlan(t *testing.T) {
	t.Parallel()
	plan, err := buildMigrationPlan(os.DirFS("./testdata").(fs.ReadDirFS))
	require.NoError(t, err)
	require.Equal(t, []MigrationDefinition{
		{
			Name:                 "001.sql",
			Content:              "create table users\n(\n    name varchar not null\n);",
			RunInsideTransaction: true,
		},
		{
			Name:                 "002.no-tx.sql",
			Content:              "CREATE INDEX CONCURRENTLY ON users (name);",
			RunInsideTransaction: false,
		},
		{
			Name:                 "003.tx.sql",
			Content:              "create table records\n(\n    field varchar\n);",
			RunInsideTransaction: true,
		},
	}, plan)
}

const checkTableExists = `SELECT EXISTS (
SELECT FROM 
	pg_tables
WHERE 
	schemaname = 'public' AND 
	tablename  = '%s'
);`

func Test_createTable(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	testDB := test.NewPostgresTest(t,
		test.WithDeleteDatabaseFunction(test.ForceDeleteDatabaseFunction),
		test.WithBaseAddress("postgres://postgres:root@localhost:1432"),
	)
	pool, err := pgxpool.New(ctx, testDB)
	require.NoError(t, err)
	// check table exists
	var exists bool
	err = pool.QueryRow(ctx, fmt.Sprintf(checkTableExists, "migrations")).Scan(&exists)
	require.NoError(t, err)
	require.False(t, exists)
	err = pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		// migrate the table
		err := createTable(ctx, conn.Conn())
		require.NoError(t, err)
		// check table exists again
		err = pool.QueryRow(ctx, fmt.Sprintf(checkTableExists, "migrations")).Scan(&exists)
		require.NoError(t, err)
		require.True(t, exists)
		// migrate again and expects no error
		return createTable(ctx, conn.Conn())
	})
	require.NoError(t, err)
}

func Test_lastMigrationApplied(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	testDB := test.NewPostgresTest(t,
		test.WithDeleteDatabaseFunction(test.ForceDeleteDatabaseFunction),
		test.WithBaseAddress("postgres://postgres:root@localhost:1432"),
	)
	pool, err := pgxpool.New(ctx, testDB)
	require.NoError(t, err)
	err = pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		// create migration table
		err := createTable(ctx, conn.Conn())
		require.NoError(t, err)
		// get the last migration applied
		migration, err := lastMigrationApplied(ctx, conn.Conn())
		require.NoError(t, err)
		require.Equal(t, "", migration)
		// insert a few records
		_, err = conn.Exec(ctx, `
INSERT INTO migrations (name) VALUES ('001.sql'), ('002.sql'), ('003.sql');
`)
		require.NoError(t, err)
		// get the last migration applied again
		migration, err = lastMigrationApplied(ctx, conn.Conn())
		require.NoError(t, err)
		require.Equal(t, "003.sql", migration)
		return nil
	})
	require.NoError(t, err)
}

func Test_insertMigrationRecord(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	testDB := test.NewPostgresTest(t,
		test.WithDeleteDatabaseFunction(test.ForceDeleteDatabaseFunction),
		test.WithBaseAddress("postgres://postgres:root@localhost:1432"),
	)
	pool, err := pgxpool.New(ctx, testDB)
	require.NoError(t, err)
	err = pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		// create migration table
		err := createTable(ctx, conn.Conn())
		require.NoError(t, err)
		// check the number of records
		var count int
		err = conn.QueryRow(ctx, `SELECT COUNT(*) FROM migrations`).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 0, count)
		// insert a record
		err = insertMigrationRecord(ctx, conn.Conn(), "test.sql")
		require.NoError(t, err)
		// check the number of records again
		err = conn.QueryRow(ctx, `SELECT COUNT(*) FROM migrations`).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)
		return nil
	})
	require.NoError(t, err)
}

func Test_applyMigration(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	testDB := test.NewPostgresTest(t,
		test.WithDeleteDatabaseFunction(test.ForceDeleteDatabaseFunction),
		test.WithBaseAddress("postgres://postgres:root@localhost:1432"),
	)
	pool, err := pgxpool.New(ctx, testDB)
	require.NoError(t, err)
	err = pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		// create migration table
		err := createTable(ctx, conn.Conn())
		require.NoError(t, err)
		// apply a migration
		err = applyMigration(ctx, conn.Conn(), MigrationDefinition{
			Name:                 "001.sql",
			Content:              "create table a(b varchar);",
			RunInsideTransaction: true,
		})
		require.NoError(t, err)
		// check if the migration was applied
		var (
			count       int
			tableExists bool
		)
		err = conn.QueryRow(ctx, `SELECT COUNT(*) FROM migrations`).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)
		err = conn.QueryRow(ctx, fmt.Sprintf(checkTableExists, "a")).Scan(&tableExists)
		require.NoError(t, err)
		require.True(t, tableExists)
		// try to execute a migration that cannot be executed inside a transaction in a transaction
		err = applyMigration(ctx, conn.Conn(), MigrationDefinition{
			Name:                 "002.sql",
			Content:              "create index concurrently on a (b);",
			RunInsideTransaction: true,
		})
		require.Error(t, err)
		require.ErrorContains(t, err, "cannot run inside a transaction block")
		// execute again outside a transaction
		err = applyMigration(ctx, conn.Conn(), MigrationDefinition{
			Name:                 "002.sql",
			Content:              "create index concurrently on a (b);",
			RunInsideTransaction: false,
		})
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)
}

func TestMigrate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	testDB := test.NewPostgresTest(t,
		test.WithDeleteDatabaseFunction(test.ForceDeleteDatabaseFunction),
		test.WithBaseAddress("postgres://postgres:root@localhost:1432"),
	)
	pool, err := pgxpool.New(ctx, testDB)
	require.NoError(t, err)
	migrationsDir := os.DirFS("./testdata").(fs.ReadDirFS)
	err = pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		// migrate
		err := Migrate(ctx, conn.Conn(), migrationsDir)
		require.NoError(t, err)
		// check the tables are created
		var (
			countFirstTime  int
			countSecondTime int
		)
		err = conn.QueryRow(ctx, `SELECT COUNT(*) FROM migrations`).Scan(&countFirstTime)
		require.NoError(t, err)
		require.NotEqual(t, 0, countFirstTime)
		// check running again causes no problems
		err = Migrate(ctx, conn.Conn(), migrationsDir)
		require.NoError(t, err)
		err = conn.QueryRow(ctx, `SELECT COUNT(*) FROM migrations`).Scan(&countSecondTime)
		require.Equal(t, countFirstTime, countSecondTime)
		return nil
	})
	require.NoError(t, err)
}
