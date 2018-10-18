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
	"regexp"
	"strconv"
	"time"

	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"github.com/spf13/cobra"
)

type cmdSnap struct {
	prefix string
	dryRun bool
	all    bool
}

func (c *cmdSnap) Command() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = "snap"
	cmd.Short = "Bulk snapshot operations"
	cmd.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println("snap")
	}
	cmd.PersistentFlags().StringVarP(&c.prefix, "prefix", "", "auto", "snapshot prefix")
	cmd.PersistentFlags().BoolVarP(&c.dryRun, "dry-run", "t", false, "dry-run don't touch a	return cmd")
	cmd.PersistentFlags().BoolVarP(&c.all, "all", "a", false, "snapshot all running containers")
	return cmd
}

type cmdSnapCreate struct {
	cmdSnap      *cmdSnap
	periodString string
	period       int
	count        int
}

func (c *cmdSnapCreate) getPeriod() int {
	p, err := parsePeriod(c.periodString)
	if err != nil {
		p = 2
	}
	return p
}

func (c *cmdSnapCreate) Command(cmdSnap *cmdSnap) *cobra.Command {
	cmd := &cobra.Command{}
	c.cmdSnap = cmdSnap
	cmd.Use = "create"
	cmd.Short = "Create automatic snapshots"
	cmd.RunE = c.Run
	cmd.Flags().StringVarP(&c.periodString, "period", "", "1s", `period in seconds, or minutes, hours, days, according to the suffix (s, m, h, d)
examples: 1h 1d`)
	cmd.Flags().IntVarP(&c.count, "count", "n", 2, "number of snapshots to keep")
	return cmd
}

func wait(op lxd.Operation, err error) error {
	if err == nil {
		return op.Wait()
	}
	return err
}

func parsePeriod(period string) (int, error) {
	re := regexp.MustCompile("([0-9]+)(.*)")
	parts := re.FindStringSubmatch(period)

	if parts != nil {
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
		suffix := parts[2]
		if suffix == "s" || suffix == "" {
			return n, nil
		} else if suffix == "m" {
			return n * 60, nil
		} else if suffix == "h" {
			return n * 3600, nil
		} else if suffix == "d" {
			return n * 3600 * 24, nil
		} else {
			return 0, errors.New("unknown suffix: " + suffix)
		}
	}
	return 0, errors.New("invalid period: " + period)
}

func (c *cmdSnapCreate) Run(cmd *cobra.Command, args []string) error {
	fmt.Println("snap create")
	server, err := GetServer()
	if err != nil {
		return err
	}
	var names []string
	if c.cmdSnap.all {
		containers, err := server.GetContainers()
		if err != nil {
			return err
		}
		for _, container := range containers {
			if container.IsActive() {
				names = append(names, container.Name)
			}
		}
	} else {
		names = args
	}

	now := time.Now()
	n := (now.Unix() / int64(c.getPeriod())) % int64(c.count)

	snapshot := api.ContainerSnapshotsPost{
		Name: c.cmdSnap.prefix + strconv.Itoa(int(n)),
	}

	if c.cmdSnap.dryRun {
		fmt.Println("snapshot name:", snapshot.Name)
	}

	for _, name := range names {
		fmt.Println(name)
		if c.cmdSnap.dryRun {
			fmt.Println(name)
		} else {
			err := wait(server.DeleteContainerSnapshot(name, snapshot.Name))
			if err != nil && "not found" != err.Error() {
				return err
			}
			err = wait(server.CreateContainerSnapshot(name, snapshot))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func init() {
	fmt.Println("snap.init")
	var snap = cmdSnap{}
	var snapCmd = snap.Command()
	rootCmd.AddCommand(snapCmd)

	var snapCreate cmdSnapCreate
	snapCreate.cmdSnap = &snap
	snapCmd.AddCommand(snapCreate.Command(&snap))
}
