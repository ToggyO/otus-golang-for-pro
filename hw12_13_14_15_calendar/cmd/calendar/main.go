package main

import (
	"context"
	"fmt"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/commands"
	"os"
)

func main() {
	defaultCmd := commands.NewAppRunCmd(Startup{})
	cmds := []commands.ICmdRunner{
		defaultCmd,
		commands.NewMigrateCmd(),
	}

	rootCmd := commands.NewRootCmd(cmds, defaultCmd)

	var subcommand = ""
	var args []string
	if len(os.Args) > 1 {
		subcommand = os.Args[1]
		args = os.Args[1:]
	}

	var err error
	if err = rootCmd.Init(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = rootCmd.Run(context.TODO(), subcommand); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
