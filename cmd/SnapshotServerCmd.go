/* SPDX-License-Identifier: Apache-2.0
*  Copyright 2018 Alex Athanasopoulos
*/
package cmd

import (
	"github.com/melato/lxdtool/op"
	"github.com/spf13/cobra"
)

func SnapshotServerCommand(s *op.Server) *cobra.Command {
	var server = &op.SnapshotServer{
		Server: s,
	}
	var command = &cobra.Command{
		RunE:  func(cmd *cobra.Command, args []string) error { return server.Run() },
		Use:   "snapshot-server",
		Short: "Handles remote requests from containers, so they can manage their own snapshots.",
		Long: `This command runs an HTTP server that listens for requests from containers to
administer their own snapshots.  It has an associated client program used by the containers.

The --profile option restricts the containers that can use this.

SECURITY WARNINGS:
- The only authentication is the client's ip address.
  That's how the server identifies the container that makes a request.
- The communication with the snapshot server is through HTTP.

ENHANCEMENTS:
- There could be a restore option, which would reboot the container
  to a snapshot.
- There be should additional authorization provided (through an access token),
  so that untrusted processes/users in the container should not be able to
  alter its snapshots.
`,
	}
	command.PersistentFlags().StringVarP(&server.Addr, "listen", "l", ":8080", "listen address")
	command.PersistentFlags().StringVarP(&server.Profile, "profile", "p", "", "profile restricting access")
	return command
}
