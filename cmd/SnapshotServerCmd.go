package cmd

import (
	"github.com/spf13/cobra"
	"melato.org/lxdtool/op"
)

func SnapshotServerCommand(s *op.Server) *cobra.Command {
	var server = &op.SnapshotServer{
		Server: s,
	}
	var command = &cobra.Command{
		Use:   "snapshot-server",
		Short: "Handles remote requests from containers, so they can manage their own snapshots.",
		RunE:  func(cmd *cobra.Command, args []string) error { return server.Run() },
	}
	command.PersistentFlags().StringVarP(&server.Addr, "listen", "l", ":8080", "listen address")
	command.PersistentFlags().StringVarP(&server.Profile, "profile", "p", "", "profile restricting access")
	return command
}
