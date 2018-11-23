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
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"gopkg.in/yaml.v2"
)

type Profile struct {
	Tool *Tool
}

func (p *Profile) List() error {
	server, err := p.Tool.GetServer()
	if err != nil {
		return err
	}
	names, err := server.GetProfileNames()
	if err != nil {
		return err
	}
	for _, name := range names {
		fmt.Println(name)
	}
	return nil
}

type ProfileExport struct {
	Tool                  *Tool
	Dir                   string
	ContainerProfilesFile string
	IncludeUsedBy         bool
}

func (c *ProfileExport) ExportProfile(server lxd.ContainerServer, name string) error {
	profile, _, err := server.GetProfile(name)
	if err != nil {
		return err
	}

	if !c.IncludeUsedBy {
		profile.UsedBy = nil
	}
	data, err := yaml.Marshal(&profile)
	if err != nil {
		return err
	}

	file := path.Join(c.Dir, name)
	fmt.Println("file", file)
	return ioutil.WriteFile(file, []byte(data), 0644)
}

func (c *ProfileExport) ExportProfiles(names []string) error {
	server, err := c.Tool.GetServer()
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

func (t *ProfileExport) ExportProfileAssociations() error {
	if t.ContainerProfilesFile == "" {
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
	f, err := os.Create(t.ContainerProfilesFile)
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
	if c.ContainerProfilesFile != "" {
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

func ImportProfile(tool *Tool, file string) error {
	server, err := tool.GetServer()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	var profile = api.Profile{}
	err = yaml.Unmarshal(data, &profile)
	if err != nil {
		return err
	}
	return server.UpdateProfile(profile.Name, profile.ProfilePut, "")
}

func ImportProfiles(tool *Tool, files []string) error {
	for _, file := range files {
		err := ImportProfile(tool, file)
		if err != nil {
			return err
		}
	}
	return nil
}
