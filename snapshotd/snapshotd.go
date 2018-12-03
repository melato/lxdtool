// SPDX-License-Identifier: Apache-2.0
// Copyright 2018 Alex Athanasopoulos
package main

import (
	"os"

	"github.com/melato/lxdtool/cmd"
	"github.com/melato/lxdtool/op"
)

func main() {
	var server op.Server
	var command = cmd.SnapshotServerCommand(&server)
	cmd.ServerFlags(command, &server)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
