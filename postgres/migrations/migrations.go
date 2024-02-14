package migrations

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/crossworth/pkgs/postgres/xpgx"
	"github.com/jackc/pgx/v5"
)

// MigrationDefinition is a migration definition.
type MigrationDefinition struct {
	Name                 string
	Content              string
	RunInsideTransaction bool
}

// migrationParts returns the parts of the given migration name.
func migrationAttributes(name string) []string {
	return strings.Split(name, ".")
}

// migrationPrefix returns the migration prefix from the given name.
func migrationPrefix(name string) string {
	return migrationAttributes(name)[0]
}

// sortMigrations sorts the migration based on the name.
func sortMigrations(migrations []string) ([]string, error) {
	var errList []error
	sort.Slice(migrations, func(i, j int) bool {
		a, err := strconv.ParseInt(migrationPrefix(migrations[i]), 10, 64)
		if err != nil {
			errList = append(errList, err)
			return false
		}
		b, err := strconv.ParseInt(migrationPrefix(migrations[j]), 10, 64)
		if err != nil {
			errList = append(errList, err)
			return false
		}
		return a < b
	})
	return migrations, errors.Join(errList...)
}

// isMigrationTx check if the migration should be executed inside a transaction given the attributes.
func isMigrationTx(attributes []string) bool {
	for _, attribute := range attributes {
		if attribute == "no-tx" {
			return false
		}
	}
	return true
}

// buildMigrationPlan builds the migration plan.
func buildMigrationPlan(folder fs.ReadDirFS) ([]MigrationDefinition, error) {
	var (
		plan  []MigrationDefinition
		files []string
	)
	entries, err := folder.ReadDir(".")
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		files = append(files, e.Name())
	}
	if files, err = sortMigrations(files); err != nil {
		return nil, err
	}
	for _, f := range files {
		file, err := folder.Open(f)
		if err != nil {
			return nil, fmt.Errorf("opening file %s: %w", f, err)
		}
		contents, err := io.ReadAll(file)
		_ = file.Close()
		if err != nil {
			return nil, fmt.Errorf("reading file %s: %w", f, err)
		}
		isTx := isMigrationTx(migrationAttributes(f))
		plan = append(plan, MigrationDefinition{
			Name:                 f,
			Content:              string(contents),
			RunInsideTransaction: isTx,
		})
	}
	return plan, nil
}

// createTable creates the migration table.
func createTable(ctx context.Context, conn xpgx.Executable) error {
	if _, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS migrations (
	id SERIAL,
	name CHARACTER VARYING NOT NULL,
	applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);
`); err != nil {
		return fmt.Errorf("creating migration table: %w", err)
	}
	return nil
}

// lastMigrationApplied returns the last migration applied.
func lastMigrationApplied(ctx context.Context, conn xpgx.Queryable) (string, error) {
	rows, _ := conn.Query(ctx, `SELECT name FROM migrations ORDER BY id DESC LIMIT 1;`)
	lastMigration, err := pgx.CollectOneRow(rows, pgx.RowTo[string])
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("querying for last migration: %w", err)
	}
	return lastMigration, nil
}

// insertMigrationRecord inserts a migration record on the migration table.
func insertMigrationRecord(ctx context.Context, conn xpgx.Executable, name string) error {
	if _, err := conn.Exec(ctx, `INSERT INTO migrations (name) VALUES ($1);`, name); err != nil {
		return fmt.Errorf("migrations: inserting migration record: %w", err)
	}
	return nil
}

// applyMigration applies a migration definition.
func applyMigration(ctx context.Context, conn xpgx.Connection, definition MigrationDefinition) error {
	if definition.RunInsideTransaction {
		if err := pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
			if _, err := tx.Exec(ctx, definition.Content); err != nil {
				return err
			}
			return insertMigrationRecord(ctx, tx, definition.Name)
		}); err != nil {
			return fmt.Errorf("executing migration: %w", err)
		}
	} else {
		if _, err := conn.Exec(ctx, definition.Content); err != nil {
			return fmt.Errorf("executing migration: %w", err)
		}
		if err := insertMigrationRecord(ctx, conn, definition.Name); err != nil {
			return fmt.Errorf("executing migration: %w", err)
		}
	}
	return nil
}

// Migrate executes the migration process on the database.
// The migration can block if huge changes are performed.
func Migrate(ctx context.Context, conn xpgx.Connection, migrationsDir fs.ReadDirFS) error {
	migrationPlan, err := buildMigrationPlan(migrationsDir)
	if err != nil {
		return fmt.Errorf("migrations: building migration plan: %w", err)
	}
	if err := createTable(ctx, conn); err != nil {
		return fmt.Errorf("migrations: creating table: %w", err)
	}
	migration, err := lastMigrationApplied(ctx, conn)
	if err != nil {
		return fmt.Errorf("migrations: querying for last migration applied: %w", err)
	}
	migrationIdx := slices.IndexFunc(migrationPlan, func(definition MigrationDefinition) bool {
		return definition.Name == migration
	})
	for _, p := range migrationPlan[migrationIdx+1:] {
		if err := applyMigration(ctx, conn, p); err != nil {
			return fmt.Errorf("migrations: applying migration %s: %w", p.Name, err)
		}
	}
	return nil
}
