package migrations

import (
	"context"

	_ "github.com/lib/pq" // lint:ignore revive
	"github.com/pressly/goose"
)

type IMigrationRunner interface {
	MigrateUp(ctx context.Context) error
	MigrateDown(ctx context.Context) error
}

type migrationRunner struct {
	dialect          string
	connectionString string
}

func NewMigrationRunner(dialect string, connectionString string) IMigrationRunner {
	return migrationRunner{
		dialect:          dialect,
		connectionString: connectionString,
	}
}

func (mr migrationRunner) MigrateUp(_ context.Context) error {
	db, err := goose.OpenDBWithDriver(mr.dialect, mr.connectionString)
	if err != nil {
		return err
	}
	return goose.Up(db, ".")
}

func (mr migrationRunner) MigrateDown(_ context.Context) error {
	db, err := goose.OpenDBWithDriver(mr.dialect, mr.connectionString)
	if err != nil {
		return err
	}
	return goose.Down(db, ".")
}
