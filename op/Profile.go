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
	"path"

	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/lxc/lxd/client"
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
	Tool *Tool
	Dir  string
}

func (c *ProfileExport) ExportProfile(server lxd.ContainerServer, name string) error {
	profile, _, err := server.GetProfile(name)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(&profile)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(c.Dir, name), []byte(data), 0644)
}

func (c *ProfileExport) Run(names []string) error {
	if c.Dir == "" {
		return errors.New("missing export dir")
	}
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
