/* SPDX-License-Identifier: Apache-2.0
*  Copyright 2018 Alex Athanasopoulos
*/
package main

import (
	"fmt"
	"os"

	"github.com/melato/lxdtool/cmd"
	"github.com/melato/lxdtool/op"
	"github.com/spf13/cobra"
)

func main() {
	var tool op.Tool
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "lxdtool",
		Short: "Miscellaneous LXD utilities for snapshots, backups, etc.",
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//	Run: func(cmd *cobra.Command, args []string) { },
	}

	cmd.ServerFlags(rootCmd, &tool.Server)
	rootCmd.AddCommand(cmd.ContainerCommand(&tool))
	rootCmd.AddCommand(cmd.ProfileCommand(&tool))
	rootCmd.AddCommand(cmd.SnapCommand(&tool))
	rootCmd.AddCommand(cmd.SnapshotServerCommand(&tool.Server))
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
