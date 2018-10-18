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
	"errors"
	"fmt"
	"path"

	"io/ioutil"

	"github.com/lxc/lxd/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "profile export, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("profile run")
	},
}

func ListProfiles() {
	server, err := GetServer()
	names, err := server.GetProfileNames()
	if err == nil {
		for _, name := range names {
			fmt.Println(name)
		}
	}
}

type cmdProfileExport struct {
	dir string
}

func (c *cmdProfileExport) Command() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = "export"
	cmd.Short = "Export profiles"
	cmd.RunE = c.Run
	cmd.PersistentFlags().StringVarP(&c.dir, "dir", "d", "", "export directory")

	return cmd
}

func (c *cmdProfileExport) ExportProfile(server lxd.ContainerServer, name string) error {
	profile, _, err := server.GetProfile(name)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(&profile)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(c.dir, name), []byte(data), 0644)
}

func (c *cmdProfileExport) Run(cmd *cobra.Command, names []string) error {
	if c.dir == "" {
		return errors.New("missing export dir")
	}
	server, err := GetServer()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		names, err = server.GetProfileNames()
		if err != nil {
			return err
		}
	}
	for _, name := range names {
		fmt.Println(name)
		err = c.ExportProfile(server, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	fmt.Println("profile.init")
	rootCmd.AddCommand(profileCmd)

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List profile names",
		Run: func(cmd *cobra.Command, args []string) {
			ListProfiles()
		},
	}

	profileCmd.AddCommand(listCmd)
	var export = cmdProfileExport{}
	profileCmd.AddCommand(export.Command())
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// profileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// profileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
