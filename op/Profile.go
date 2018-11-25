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
	Tool          *Tool
	File          string
	Dir           string
	All           bool
	IncludeUsedBy bool
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

func ProfileExists(server lxd.ContainerServer, name string) (bool, error) {
	names, err := server.GetProfileNames()
	if err != nil {
		return false, err
	}
	for _, profile := range names {
		if profile == name {
			return true, nil
		}
	}
	return false, nil
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
	name := profile.Name
	exists, err := ProfileExists(server, name)
	if err != nil {
		return err
	}
	if exists {
		return server.UpdateProfile(profile.Name, profile.ProfilePut, "")
	} else {
		var post api.ProfilesPost
		post.ProfilePut = profile.ProfilePut
		post.Name = name
		return server.CreateProfile(post)
	}
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
