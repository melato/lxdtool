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
	"melato.org/lxdtool/op"
)

func ContainerCommand1(tool *op.Tool) *cobra.Command {
	containerCmd := &cobra.Command{}
	containerCmd.Use = "container"
	containerCmd.Short = "List containers"
	containerCmd.Run = func(cmd *cobra.Command, args []string) {
		tool.ListContainers(args)
	}

	return containerCmd
}

func ListCommand(tool *op.Tool) *cobra.Command {
	listCmd := &cobra.Command{}
	listCmd.Use = "list"
	listCmd.Run = func(cmd *cobra.Command, args []string) {
		tool.ListContainers(args)
	}
	return listCmd
}

func ContainerCommand(tool *op.Tool) *cobra.Command {
	containerCmd := &cobra.Command{}
	containerCmd.Use = "container"
	//rootCmd.AddCommand(containerCmd)

	listCmd := ListCommand(tool)
	containerCmd.AddCommand(listCmd)
	//rootCmd.AddCommand(listCmd)

	profilesCmd := &cobra.Command{}
	profilesCmd.Use = "profiles"
	profilesCmd.Run = func(cmd *cobra.Command, args []string) {
		tool.ListContainerProfiles(args)
	}
	containerCmd.AddCommand(profilesCmd)

	addressesCmd := &cobra.Command{}
	addressesCmd.Use = "addresses"
	addressesCmd.Run = func(cmd *cobra.Command, args []string) {
		tool.ListContainerAddresses(args)
	}
	containerCmd.AddCommand(addressesCmd)

	ip4Cmd := &cobra.Command{}
	ip4Cmd.Use = "ip4"
	ip4Cmd.Run = func(cmd *cobra.Command, args []string) {
		tool.ListContainerAddressesIP4(args)
	}
	containerCmd.AddCommand(ip4Cmd)

	pidCmd := &cobra.Command{}
	pidCmd.Use = "pid"
	pidCmd.Run = func(cmd *cobra.Command, args []string) {
		tool.ListContainerPids(args)
	}
	containerCmd.AddCommand(pidCmd)

	findCmd := &cobra.Command{}
	findCmd.Use = "find"
	findCmd.Run = func(cmd *cobra.Command, args []string) {
		tool.FindPids(args)
	}
	containerCmd.AddCommand(findCmd)
	return containerCmd
}
