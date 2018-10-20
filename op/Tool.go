package op

import (
	"fmt"

	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

type Tool struct {
	SocketPath string
	server     lxd.ContainerServer
	All        bool
}

func (c *Tool) GetServer() (lxd.ContainerServer, error) {
	if c.server == nil {
		// Connect to LXD over the Unix socket
		var err error
		c.server, err = lxd.ConnectLXDUnix(c.SocketPath, nil)
		if err != nil {
			return nil, err
		}
	}
	return c.server, nil
}

func (c *Tool) GetContainerNames(args []string) ([]string, error) {
	if c.All {
		var names []string
		server, err := c.GetServer()
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

func (c *Tool) ListContainers(args []string) error {
	names, err := c.GetContainerNames(args)
	if err != nil {
		return err
	}
	for _, name := range names {
		fmt.Println(name)
	}
	return nil
}

func (t *Tool) ListContainerProfiles(args []string) error {
	containers, err := t.GetContainerNames(args)
	if err != nil {
		return err
	}
	server, err := t.GetServer()
	if err != nil {
		return err
	}
	for _, name := range containers {
		c, _, err := server.GetContainer(name)
		if err != nil {
			return err
		}
		for _, profile := range c.Profiles {
			fmt.Println(name, profile)
		}
	}
	return nil
}

func (t *Tool) ListContainerAddressesF(args []string, includeNetworkName bool, f func(api.ContainerStateNetworkAddress) interface{}) error {
	containers, err := t.GetContainerNames(args)
	if err != nil {
		return err
	}
	server, err := t.GetServer()
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
func (t *Tool) ListContainerAddresses(args []string) error {
	return t.ListContainerAddressesF(args, true, func(address api.ContainerStateNetworkAddress) interface{} {
		return address
	})
}

func (t *Tool) ListContainerAddressesIP4(args []string) error {
	return t.ListContainerAddressesF(args, false, func(address api.ContainerStateNetworkAddress) interface{} {
		if address.Family == "inet" && address.Scope == "global" {
			return address.Address
		}
		return nil
	})
}
