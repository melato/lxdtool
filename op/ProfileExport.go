/* SPDX-License-Identifier: Apache-2.0
*  Copyright 2018 Alex Athanasopoulos
*/
package op

import (
	"encoding/csv"
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/lxc/lxd/client"
	"gopkg.in/yaml.v2"
)

type ProfileExport struct {
	Tool *Tool
	File string
	Dir  string
	All  bool
}

func (c *ProfileExport) ExportProfile(server lxd.ContainerServer, name string) error {
	profile, _, err := server.GetProfile(name)
	if err != nil {
		return err
	}

	// remove the UsedBy info
	profile.UsedBy = nil
	data, err := yaml.Marshal(&profile)
	if err != nil {
		return err
	}

	file := path.Join(c.Dir, name)
	return ioutil.WriteFile(file, []byte(data), 0644)
}

func (c *ProfileExport) ExportProfiles(names []string) error {
	server, err := c.Tool.GetServer()
	if err != nil {
		return err
	}
	if len(names) == 0 && c.All {
		names, err = server.GetProfileNames()
		if err != nil {
			return err
		}
	}
	for _, name := range names {
		err = c.ExportProfile(server, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *ProfileExport) ExportProfileAssociations() error {
	if t.File == "" {
		return errors.New("missing export file")
	}
	server, err := t.Tool.GetServer()
	if err != nil {
		return err
	}
	containers, err := server.GetContainers()
	if err != nil {
		return err
	}
	f, err := os.Create(t.File)
	if err != nil {
		return err
	}
	defer f.Close()
	cs := csv.NewWriter(f)
	defer cs.Flush()
	var row = []string{"container", "profile"}
	cs.Write(row)
	for _, container := range containers {
		row[0] = container.Name
		for _, profile := range container.Profiles {
			row[1] = profile
			cs.Write(row)
		}
	}
	return nil
}

func (c *ProfileExport) Run(names []string) error {
	if c.File != "" {
		err := c.ExportProfileAssociations()
		if err != nil {
			return err
		}
	}
	err := c.ExportProfiles(names)
	if err != nil {
		return err
	}
	return nil
}
