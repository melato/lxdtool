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
	Exclude    []string
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

func StringSliceDiff(ar []string, exclude []string) []string {
	if exclude == nil {
		return ar
	}
	var xmap = make(map[string]bool)
	for _, s := range exclude {
		xmap[s] = true
	}
	var result []string
	for _, s := range ar {
		if !xmap[s] {
			result = append(result, s)
		}
	}
	return result
}

func (c *Tool) GetContainerNames(args []string) ([]string, error) {
	var names []string
	if c.All {
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
	} else {
		names = args
	}
	return StringSliceDiff(names, c.Exclude), nil
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

func (t *Tool) ListContainerPids(args []string) error {
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
		fmt.Println(name, state.Pid)
	}
	return nil
}
