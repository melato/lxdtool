/* Copyright 2018 Alex Athanasopoulos
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*   http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
 */
package op

import (
	"io/ioutil"

	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"gopkg.in/yaml.v2"
)

type ProfileImport struct {
	Tool   *Tool
	Update bool
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

func (t *ProfileImport) ImportProfile(file string) error {
	server, err := t.Tool.GetServer()
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

func (t *ProfileImport) ImportProfiles(files []string) error {
	for _, file := range files {
		err := t.ImportProfile(file)
		if err != nil {
			return err
		}
	}
	return nil
}
