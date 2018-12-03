// SPDX-License-Identifier: Apache-2.0
// Copyright 2018 Alex Athanasopoulos
package cmd

import (
	"fmt"

	"github.com/melato/lxdtool/op"
	"github.com/spf13/cobra"
)

func profileExportCommand(t *op.ProfileExport) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = "export [flags] [profile] ..."
	cmd.Short = "Export profiles"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return t.Run(args)
	}
	cmd.PersistentFlags().StringVarP(&t.Dir, "dir", "d", "", "export directory")
	cmd.PersistentFlags().BoolVarP(&t.All, "all", "a", false, "export all profiles")
	cmd.PersistentFlags().StringVarP(&t.File, "file", "f", "", "container-profile csv file")
	return cmd
}

func profileImportCommand(t *op.ProfileImport) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = "import [flags] [profile-file] ..."
	cmd.Short = "Create or Update profiles"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return t.ImportProfiles(args)
	}
	cmd.PersistentFlags().BoolVarP(&t.Update, "update", "u", false, "update existing profiles")
	return cmd
}

func ProfileCommand(tool *op.Tool) *cobra.Command {
	var opProfile = op.Profile{tool}
	var profileCmd = &cobra.Command{
		Use:   "profile",
		Short: "profile export, etc.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("profile run")
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List profile names",
		Run: func(cmd *cobra.Command, args []string) {
			opProfile.List()
		},
	}
	profileCmd.AddCommand(listCmd)
	var opExport = &op.ProfileExport{}
	opExport.Tool = tool
	profileCmd.AddCommand(profileExportCommand(opExport))
	var opImport = &op.ProfileImport{}
	opImport.Tool = tool
	profileCmd.AddCommand(profileImportCommand(opImport))
	return profileCmd
}
