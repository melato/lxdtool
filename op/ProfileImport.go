// SPDX-License-Identifier: Apache-2.0
// Copyright 2018 Alex Athanasopoulos
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
