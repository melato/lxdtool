/* Copyright 2018 Alex Athanasopoulos
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*   http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
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
	rootCmd.AddCommand(cmd.TestCommand())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
