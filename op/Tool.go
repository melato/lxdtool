package op

import (
	"github.com/lxc/lxd/client"
)

type Tool struct {
	Server Server
}

func (t *Tool) GetServer() (lxd.ContainerServer, error) {
	return t.Server.GetServer()
}
