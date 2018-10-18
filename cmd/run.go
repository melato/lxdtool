package cmd

import (
	"github.com/spf13/cobra"
)

type RunOp func(args []string) error

func WrapRun(f RunOp) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return f(args)
	}
}
