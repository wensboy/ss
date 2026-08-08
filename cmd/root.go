package cmd

import (
	"context"
	"os"

	"github.com/wensboy/ss/config"

	"github.com/urfave/cli/v3"
)

func Execute() {
	rootCmd := config.InitCommand("spec/command.json")

	mountCmd(rootCmd)

	if err := rootCmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}

func mountCmd(cmd *cli.Command) {
	cmd.Action = func(c context.Context, cliCmd *cli.Command) error {
		return cli.Exit("type 'ss help' for more information\n", 1)
	}
}
