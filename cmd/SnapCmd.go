/* SPDX-License-Identifier: Apache-2.0
*  Copyright 2018 Alex Athanasopoulos
 */
package cmd

import (
	"github.com/melato/lxdtool/op"
	"github.com/spf13/cobra"
)

func CreateCommand(c *op.SnapCreate, opSnap *op.Snap) *cobra.Command {
	cmd := &cobra.Command{}
	c.Snap = opSnap
	cmd.Use = "create [flags] [containers]"
	cmd.Short = "Create automatic snapshots"
	cmd.Long = `Creates snapshots for specified containers, using a rotating naming scheme.
The snapshot names are determined by the appending a numeric suffix to the {prefix},
representing the current time {period} modulo {count},
so that if the command is executed periodically every {period},
there would be {count} different snapshots.
Any previous snapshot with the same name is deleted.
The command is meant to be run periodically, at the same frequency as specified in the {period}.`
	cmd.Example = `lxdtool snap create my-container
lxdtool snap create -a --period 1h --count 24 --prefix auto_hour
lxdtool snap create -a --period 1d --count 7 --prefix auto_day
lxdtool snap create -a --period 7d --count 4 --prefix auto_week`
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return c.Run(args)
	}
	cmd.Flags().StringVarP(&c.PeriodString, "period", "", "1s", `period in seconds, or minutes, hours, days, according to the suffix (s, m, h, d)
examples: 1h 1d`)
	cmd.Flags().IntVarP(&c.Count, "count", "n", 0, "number of snapshots to keep.  0 means use no prefix")
	return cmd
}

func DeleteCommand(c *op.SnapDelete, opSnap *op.Snap) *cobra.Command {
	cmd := &cobra.Command{}
	c.Snap = opSnap
	cmd.Use = "delete"
	cmd.Short = "Delete automatic snapshots"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return c.Run(args)
	}
	return cmd
}

func SnapCommand(tool *op.Tool) *cobra.Command {
	var snap = &op.Snap{}
	snap.Tool = tool
	snapCmd := &cobra.Command{}
	snapCmd.Use = "snap"
	snapCmd.Short = "Bulk snapshot operations"
	ContainerFlags(snapCmd, &snap.ContainerOptions)
	snapCmd.PersistentFlags().StringVar(&snap.Prefix, "name", "auto", "snapshot name or prefix")
	snapCmd.PersistentFlags().BoolVarP(&snap.DryRun, "dry-run", "t", false, "dry-run don't touch a	return cmd")

	var snapCreate = &op.SnapCreate{}
	snapCreate.Snap = snap
	snapCmd.AddCommand(CreateCommand(snapCreate, snap))

	var snapDelete = &op.SnapDelete{}
	snapDelete.Snap = snap
	snapCmd.AddCommand(DeleteCommand(snapDelete, snap))
	return snapCmd
}
