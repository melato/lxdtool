/* SPDX-License-Identifier: Apache-2.0
*  Copyright 2018 Alex Athanasopoulos
*/
package cmd

import (
	"github.com/melato/lxdtool/op"
	"github.com/spf13/cobra"
)

func ServerFlags(cmd *cobra.Command, server *op.Server) {
	cmd.PersistentFlags().StringVarP(&server.Socket, "socket", "s", "/var/snap/lxd/common/lxd/unix.socket", "path to unix socket")
	cmd.PersistentFlags().StringVar(&server.Remote, "remote", "", "LXD remote")
	cmd.PersistentFlags().StringVarP(&server.ConfigFile, "config", "c", "${HOME}/snap/lxd/current/.config/lxc/config.yml", "path to config.yml")
}
