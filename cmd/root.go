package cmd

import (
	"context"
	"fmt"
	"os"

	"example.com/s_xiewenjun/opt/config"

	"github.com/urfave/cli/v3"
)

func Execute() {
	rootCmd := config.InitCommand("spec/command.json")

	mountRoot(rootCmd)

	if err := rootCmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}

func mountRoot(root *cli.Command) {
	root.Action = func(c context.Context, cmd *cli.Command) error {
		fmt.Fprintf(os.Stderr, "nothing to do~\n")
		return nil
	}
}
