package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/application"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/configuration"
)

const (
	appRunCmdName        = "run"
	appRunCmdDescription = "Build and run application."
	appRunCmdUsage       = `Usage:
        <rootcmd> [--root_options]... [command] [--options]...
Commands:
	run                         Build and run application.

Root options:
	-h, --help                  Show command line help.

Options:
	-c, --config                Path to configuration file.`
)

var ErrInvalidAppRunCmdArgs = errors.New("invalid arguments provided. Call `<rootcmd> --help` to list possible actions")

type appRunCmd struct {
	fs *flag.FlagSet

	description string
	configPath  string
	help        bool

	startup            application.IStartup
	applicationBuilder application.IApplicationBuilder
}

func NewAppRunCmd(startup application.IStartup) ICmdRunner {
	fs := flag.NewFlagSet(appRunCmdName, flag.ExitOnError)

	ac := &appRunCmd{
		fs:          fs,
		description: appRunCmdDescription,
		startup:     startup,
	}

	ac.fs.BoolVar(&ac.help, cmdHelp, false, "show command line help")
	ac.fs.BoolVar(&ac.help, cmdHelpShort, false, "show command line help")
	ac.fs.StringVar(&ac.configPath, "config", defaultCfgPath, cfgFlagUsage)
	ac.fs.StringVar(&ac.configPath, "c", defaultCfgPath, cfgFlagUsage)

	return ac
}

func (ac *appRunCmd) Name() string {
	return ac.fs.Name()
}

func (ac *appRunCmd) Description() string {
	return ac.description
}

func (ac *appRunCmd) Init(args []string) error {
	if err := ac.fs.Parse(args); err != nil {
		return err
	}

	if ac.help && len(args) > 1 {
		return ErrInvalidAppRunCmdArgs
	}

	config := configuration.NewConfiguration(ac.configPath)
	ac.applicationBuilder = application.NewApplicationBuilder(config, ac.startup)

	return nil
}

func (ac *appRunCmd) Run(_ context.Context) error {
	if ac.help {
		fmt.Println(appRunCmdUsage)
		return nil
	}

	app, err := ac.applicationBuilder.Build()
	if err != nil {
		return err
	}

	app.Run()
	return nil
}
