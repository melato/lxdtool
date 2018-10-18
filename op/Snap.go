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

package op

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

type Snap struct {
	Tool   *Tool
	Prefix string
	DryRun bool
	All    bool
}

func (c *Snap) GetContainerNames(args []string) ([]string, error) {
	if c.All {
		var names []string
		server, err := c.Tool.GetServer()
		if err != nil {
			return nil, err
		}
		containers, err := server.GetContainers()
		if err != nil {
			return nil, err
		}
		for _, container := range containers {
			if container.IsActive() {
				names = append(names, container.Name)
			}
		}
		return names, nil
	} else {
		return args, nil
	}
}

type SnapCreate struct {
	Snap         *Snap
	PeriodString string
	Period       int
	Count        int
}

func (c *SnapCreate) getPeriod() int {
	p, err := parsePeriod(c.PeriodString)
	if err != nil {
		p = 2
	}
	return p
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

func (c *SnapCreate) Run(args []string) error {
	fmt.Println("snap create")
	server, err := c.Snap.Tool.GetServer()
	if err != nil {
		return err
	}
	names, err := c.Snap.GetContainerNames(args)
	if err != nil {
		return err
	}

	now := time.Now()
	snapshotName := c.Snap.Prefix
	if c.Count > 0 {
		n := (now.Unix() / int64(c.getPeriod())) % int64(c.Count)
		snapshotName = c.Snap.Prefix + strconv.Itoa(int(n))
	}

	snapshot := api.ContainerSnapshotsPost{
		Name: snapshotName,
	}

	if c.Snap.DryRun {
		fmt.Println("snapshot name:", snapshot.Name)
	}

	for _, name := range names {
		fmt.Println(name)
		if c.Snap.DryRun {
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

type SnapDelete struct {
	Snap *Snap
}

func (c *SnapDelete) Run(args []string) error {
	server, err := c.Snap.Tool.GetServer()
	if err != nil {
		return err
	}
	names, err := c.Snap.GetContainerNames(args)
	if err != nil {
		return err
	}
	for _, name := range names {
		snapshots, err := server.GetContainerSnapshotNames(name)
		if err != nil {
			return err
		}
		for _, snap := range snapshots {
			if strings.HasPrefix(snap, c.Snap.Prefix) {
				fmt.Println(name + "/" + snap)
				if !c.Snap.DryRun {
					err := wait(server.DeleteContainerSnapshot(name, snap))
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}