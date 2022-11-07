package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/configuration"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/migrations"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/shared"
)

const (
	migrateCmdName        = "migrate"
	migrateCmdDescription = "Manage database migrations."
	migrateCmdUsage       = `Usage:
        <rootcmd> migrate [-action] [--options]...
Actions:
	-up                         Run all pending migrations.
	-down                       Undo all migrations.
Options:
	-c, --config                Path to configuration file.`
)

var ErrInvalidMigrateAction = errors.
	New("invalid action flag provided. Call `<rootcmd> migrate --help` to list possible actions")

type migrateCmd struct {
	fs *flag.FlagSet

	description string
	up          bool
	down        bool
	configPath  string
	help        bool

	migrationRunner migrations.IMigrationRunner
}

func NewMigrateCmd() ICmdRunner {
	fs := flag.NewFlagSet(migrateCmdName, flag.ExitOnError)
	fs.Usage = func() {
		fmt.Print(migrateCmdUsage)
	}

	mc := &migrateCmd{
		description: migrateCmdDescription,
		fs:          fs,
	}

	mc.fs.BoolVar(&mc.down, "down", false, "undo all migrations")
	mc.fs.BoolVar(&mc.up, "up", false, "run all pending migrations")
	mc.fs.BoolVar(&mc.help, cmdHelp, false, "show command line help")
	mc.fs.BoolVar(&mc.help, cmdHelpShort, false, "show command line help")
	mc.fs.StringVar(&mc.configPath, "config", defaultCfgPath, cfgFlagUsage)
	mc.fs.StringVar(&mc.configPath, "c", defaultCfgPath, cfgFlagUsage)

	return mc
}

func (mc *migrateCmd) Name() string {
	return mc.fs.Name()
}

func (mc *migrateCmd) Description() string {
	return mc.description
}

func (mc *migrateCmd) Init(args []string) error {
	if err := mc.fs.Parse(args); err != nil {
		return err
	}

	if (mc.help && len(args) > 1) || (mc.up && mc.down) {
		return ErrInvalidMigrateAction
	}

	config := configuration.NewConfiguration(mc.configPath)
	mc.migrationRunner = migrations.NewMigrationRunner(config.Storage.Dialect,
		shared.CreatePgConnectionString(config.Storage))

	return nil
}

func (mc *migrateCmd) Run(ctx context.Context) error {
	if mc.help {
		fmt.Println(migrateCmdUsage)
		return nil
	}

	if mc.up {
		return mc.migrationRunner.MigrateUp(ctx)
	}
	return mc.migrationRunner.MigrateDown(ctx)
}
