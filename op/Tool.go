package op

import (
	"github.com/lxc/lxd/client"
)

type Tool struct {
	SocketPath string
	server     lxd.ContainerServer
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
