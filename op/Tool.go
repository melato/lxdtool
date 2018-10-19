package op

import (
	"fmt"

	"github.com/lxc/lxd/client"
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
