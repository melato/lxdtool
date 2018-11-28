package main

import (
	"os"

	"melato.org/lxdtool/cmd"
	"melato.org/lxdtool/op"
)

func main() {
	var server op.Server
	var command = cmd.SnapshotServerCommand(&server)
	cmd.ServerFlags(command, &server)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
