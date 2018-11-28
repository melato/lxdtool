package op

import (
	"fmt"

	"github.com/lxc/lxd/shared/api"
)

type ContainerOptions struct {
	All     bool
	Exclude []string
	ProcDir string
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

func (t *ContainerOps) ListContainerProfiles(args []string) error {
	containers, err := t.GetContainerNames(args)
	if err != nil {
		return err
	}
	server, err := t.Tool.GetServer()
	if err != nil {
		return err
	}
	for _, name := range containers {
		t, _, err := server.GetContainer(name)
		if err != nil {
			return err
		}
		for _, profile := range t.Profiles {
			fmt.Println(name, profile)
		}
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
