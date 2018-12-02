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
package cmd

import (
	"github.com/melato/lxdtool/op"
	"github.com/spf13/cobra"
)

func ContainerFlags(cmd *cobra.Command, c *op.ContainerOptions) {
	cmd.PersistentFlags().StringVar(&c.ProcDir, "proc", "/proc", "server /proc dir")
	cmd.PersistentFlags().BoolVarP(&c.All, "all", "a", false, "use all containers")
	cmd.PersistentFlags().BoolVarP(&c.Running, "running", "r", false, "use only running containers")
	cmd.PersistentFlags().StringVarP(&c.Profile, "profile", "p", "", "use containers that have a profile")
	cmd.PersistentFlags().StringSliceVarP(&c.Exclude, "exclude", "x", nil, "exclude containers")
}

func ListCommand(c *op.ContainerOps) *cobra.Command {
	listCmd := &cobra.Command{}
	listCmd.Use = "list"
	listCmd.Run = func(cmd *cobra.Command, args []string) {
		c.ListContainers(args)
	}
	return listCmd
}

func ContainerCommand(tool *op.Tool) *cobra.Command {
	var c = &op.ContainerOps{
		Tool: tool,
	}
	containerCmd := &cobra.Command{}
	containerCmd.Use = "container"
	ContainerFlags(containerCmd, &c.ContainerOptions)

	listCmd := ListCommand(c)
	containerCmd.AddCommand(listCmd)

	profilesCmd := &cobra.Command{}
	profilesCmd.Use = "profiles"
	profilesCmd.Run = func(cmd *cobra.Command, args []string) {
		c.ListContainerProfiles(args)
	}
	containerCmd.AddCommand(profilesCmd)

	addressesCmd := &cobra.Command{}
	addressesCmd.Use = "addresses"
	addressesCmd.Run = func(cmd *cobra.Command, args []string) {
		c.ListContainerAddresses(args)
	}
	containerCmd.AddCommand(addressesCmd)

	ip4Cmd := &cobra.Command{}
	ip4Cmd.Use = "ip4"
	ip4Cmd.Run = func(cmd *cobra.Command, args []string) {
		c.ListContainerAddressesIP4(args)
	}
	containerCmd.AddCommand(ip4Cmd)

	pidCmd := &cobra.Command{}
	pidCmd.Use = "pid"
	pidCmd.Run = func(cmd *cobra.Command, args []string) {
		c.ListContainerPids(args)
	}
	containerCmd.AddCommand(pidCmd)

	findCmd := &cobra.Command{}
	findCmd.Use = "find"
	findCmd.Run = func(cmd *cobra.Command, args []string) {
		c.Tool.Server.FindPids(c.ContainerOptions.ProcDir, args)
	}
	findCmd.Aliases = []string{"ps"}
	containerCmd.AddCommand(findCmd)
	return containerCmd
}
