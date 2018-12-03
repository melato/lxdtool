/* SPDX-License-Identifier: Apache-2.0
*  Copyright 2018 Alex Athanasopoulos
 */
package op

import (
	"errors"
	"fmt"

	"github.com/lxc/lxd/shared/api"
)

type ContainerOptions struct {
	All     bool
	Running bool
	Profile string
	Exclude []string

	ProcDir string // used by find cmd
}

type ContainerOps struct {
	Tool *Tool
	ContainerOptions
}

func (t *ContainerOps) GetContainerNames(args []string) ([]string, error) {
	return t.Tool.Server.GetContainerNames(&t.ContainerOptions, args)
}

func (t *ContainerOps) ListContainers(args []string) error {
	names, err := t.GetContainerNames(args)
	if err != nil {
		return err
	}
	for _, name := range names {
		fmt.Println(name)
	}
	return nil
}

func GetContainerName(args []string) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	}
	return "", errors.New("please specify one container")
}

func (t *ContainerOps) ListContainerProfiles(args []string) error {
	name, err := GetContainerName(args)
	if err != nil {
		return err
	}
	server, err := t.Tool.GetServer()
	if err != nil {
		return err
	}
	c, _, err := server.GetContainer(name)
	if err != nil {
		return err
	}
	for _, profile := range c.Profiles {
		fmt.Println(profile)
	}
	return nil
}

func (t *ContainerOps) ListContainerAddressesF(args []string, includeNetworkName bool, f func(api.ContainerStateNetworkAddress) interface{}) error {
	containers, err := t.GetContainerNames(args)
	if err != nil {
		return err
	}
	server, err := t.Tool.GetServer()
	if err != nil {
		return err
	}
	for _, name := range containers {
		state, _, err := server.GetContainerState(name)
		if err != nil {
			return err
		}
		for networkName, network := range state.Network {
			for _, address := range network.Addresses {
				x := f(address)
				if x != nil {
					if includeNetworkName {
						fmt.Println(name, networkName, x)
					} else {
						fmt.Println(x, name)
					}
				}
			}
		}
	}
	return nil
}
func (t *ContainerOps) ListContainerAddresses(args []string) error {
	return t.ListContainerAddressesF(args, true, func(address api.ContainerStateNetworkAddress) interface{} {
		return address
	})
}

func (t *ContainerOps) ListContainerAddressesIP4(args []string) error {
	return t.ListContainerAddressesF(args, false, func(address api.ContainerStateNetworkAddress) interface{} {
		if address.Family == "inet" && address.Scope == "global" {
			return address.Address
		}
		return nil
	})
}

func (t *ContainerOps) ListContainerPids(args []string) error {
	containers, err := t.GetContainerNames(args)
	if err != nil {
		return err
	}
	server, err := t.Tool.GetServer()
	if err != nil {
		return err
	}
	for _, name := range containers {
		state, _, err := server.GetContainerState(name)
		if err != nil {
			return err
		}
		fmt.Println(name, state.Pid)
	}
	return nil
}
