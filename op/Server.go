// SPDX-License-Identifier: Apache-2.0
// Copyright 2018 Alex Athanasopoulos
package op

import (
	"fmt"
	"os"

	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/lxc/config"
	"github.com/lxc/lxd/shared/api"
)

type Server struct {
	Socket     string
	Remote     string
	ConfigFile string
	server     lxd.ContainerServer
}

func (t *Server) GetServer() (lxd.ContainerServer, error) {
	if t.server == nil {
		var err error
		if t.Remote != "" && t.ConfigFile != "" {
			fmt.Println("using remote: ", t.Remote)
			configFile := os.ExpandEnv(t.ConfigFile)
			fmt.Println("ConfigFile: ", configFile)
			conf, err := config.LoadConfig(configFile)
			if err != nil {
				return nil, err
			}
			t.server, err = conf.GetContainerServer(t.Remote)
			if err != nil {
				return nil, err
			}
		} else {
			// Connect to LXD over the Unix socket
			t.server, err = lxd.ConnectLXDUnix(t.Socket, nil)
			if err != nil {
				fmt.Println(t.Socket, err)
				return nil, err
			}
		}
	}
	return t.server, nil
}

func HasProfile(container *api.Container, profile string) bool {
	for _, p := range container.Profiles {
		if p == profile {
			return true
		}
	}
	return false
}

func (t *Server) isSelected(opt *ContainerOptions, container *api.Container) bool {
	if opt.Running && !container.IsActive() {
		return false
	}
	if opt.Profile != "" {
		return HasProfile(container, opt.Profile)
	}
	if opt.All {
		return true
	}
	return false
}

func (t *Server) GetContainerNames(opt *ContainerOptions, args []string) ([]string, error) {
	var names []string
	if opt.All || opt.Running || (opt.Profile != "") {
		server, err := t.GetServer()
		if err != nil {
			return nil, err
		}
		containers, err := server.GetContainers()
		if err != nil {
			return nil, err
		}
		for _, container := range containers {
			if t.isSelected(opt, &container) {
				names = append(names, container.Name)
			}
		}
	} else {
		names = args
	}
	return StringSliceDiff(names, opt.Exclude), nil
}
