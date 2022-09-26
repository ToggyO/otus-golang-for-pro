package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/shared"
)

const (
	rootCmdName         = "root"
	rootCmdDescription  = "root description"
	rootCmdVersion      = "version"
	rootCmdVersionShort = "v"
)

var ErrInvalidRootOption = errors.New("invalid option flag provided. Call `<rootcmd> --help` to list possible options")

type RootCmd struct {
	fs *flag.FlagSet

	runners         []ICmdRunner
	defaultCmd      ICmdRunner
	forceRunDefault bool

	help    bool
	version bool
	usage   string
}

func NewRootCmd(runners []ICmdRunner, defaultCmd ICmdRunner) *RootCmd {
	if defaultCmd == nil && len(runners) > 0 {
		defaultCmd = runners[0]
	}

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	usage := createUsage(runners)
	rc := &RootCmd{
		fs:         fs,
		runners:    runners,
		defaultCmd: defaultCmd,
		usage:      usage,
	}

	rc.fs.BoolVar(&rc.help, cmdHelp, false, "show command line help")
	rc.fs.BoolVar(&rc.help, cmdHelpShort, false, "show command line help")
	rc.fs.BoolVar(&rc.version, rootCmdVersion, false, "show application version in use")
	rc.fs.BoolVar(&rc.version, rootCmdVersionShort, false, "show application version in use")

	return rc
}

func (rc *RootCmd) Name() string {
	return rootCmdName
}

func (rc *RootCmd) Description() string {
	return rootCmdDescription
}

func (rc *RootCmd) Init(args []string) error {
	rc.fs.SetOutput(io.Discard)
	err := rc.fs.Parse(args)
	rc.fs.SetOutput(os.Stdout)

	if err != nil && !rc.isSubCommand(rc.fs.Arg(0)) {
		rc.forceRunDefault = true
		return nil
	}

	if rc.help && rc.version {
		return ErrInvalidRootOption
	}

	return nil
}

func (rc *RootCmd) Run(ctx context.Context, subcommand string) error {
	if rc.help {
		fmt.Println(rc.usage)
		return nil
	}

	if rc.version {
		shared.PrintVersion()
		return nil
	}

	if rc.forceRunDefault || subcommand == "" {
		return rc.runSubcommand(ctx, rc.defaultCmd, os.Args[1:]...)
	}

	for _, cmd := range rc.runners {
		if cmd.Name() == subcommand {
			return rc.runSubcommand(ctx, cmd, os.Args[2:]...)
		}
	}

	return nil
}

func (rc *RootCmd) runSubcommand(ctx context.Context, cmd ICmdRunner, args ...string) error {
	err := cmd.Init(args)
	if err != nil {
		return err
	}
	err = cmd.Run(ctx)
	if err != nil {
		return err
	}
	return nil
}

func createUsage(runners []ICmdRunner) string {
	spacesFromLeftToCmd := 8
	spacesFromLeftToDescription := 36
	builder := strings.Builder{}

	builder.WriteString(`Usage:
	<rootcmd> [--options]... [command]
Commands:
`)

	for _, cmd := range runners {
		name := cmd.Name()
		description := cmd.Description()
		totalSpacesFromLeft := spacesFromLeftToCmd + len(name)
		spacesBetweenNameAndDescr := spacesFromLeftToDescription - totalSpacesFromLeft + len(description)
		builder.WriteString(fmt.Sprintf("%*s%*s\n", totalSpacesFromLeft, name, spacesBetweenNameAndDescr, description))
	}

	builder.WriteString(`Options:
	-h, --help                  Show command line help.
	-v, --version		    Show application version in use.
	-c, --config                Path to configuration file.`)

	return builder.String()
}

func (rc *RootCmd) isSubCommand(arg string) bool {
	for _, runner := range rc.runners {
		if runner.Name() == arg {
			return true
		}
	}

	return false
}
