// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
)

func ContainerCommand() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = "container"
	cmd.Short = "List containers"
	cmd.Run = func(cmd *cobra.Command, args []string) {
		tool.ListContainers(args)
	}
	return cmd
}

func init() {
	containerCmd := &cobra.Command{}
	containerCmd.Use = "container"
	rootCmd.AddCommand(containerCmd)

	listCmd := &cobra.Command{}
	listCmd.Use = "container"
	listCmd.Run = func(cmd *cobra.Command, args []string) {
		tool.ListContainers(args)
	}
	containerCmd.AddCommand(listCmd)

	profilesCmd := &cobra.Command{}
	profilesCmd.Use = "profiles"
	profilesCmd.Run = func(cmd *cobra.Command, args []string) {
		tool.ListContainerProfiles(args)
	}
	containerCmd.AddCommand(profilesCmd)
}
