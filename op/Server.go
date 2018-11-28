package op

import (
	"fmt"
	"os"
	"path"

	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/lxc/config"
)

type Server struct {
	Socket    string
	Remote    string
	ConfigDir string
	server    lxd.ContainerServer
}

func (t *Server) GetServer() (lxd.ContainerServer, error) {
	if t.server == nil {
		var err error
		if t.Remote != "" && t.ConfigDir != "" {
			fmt.Println("using remote: ", t.Remote)
			fmt.Println("ConfigDir: ", t.ConfigDir)
			confPath := os.ExpandEnv(path.Join(t.ConfigDir, "config.yml"))
			conf, err := config.LoadConfig(confPath)
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
				return nil, err
			}
		}
	}
	return t.server, nil
}

func (t *Server) GetContainerNames(opt *ContainerOptions, args []string) ([]string, error) {
	var names []string
	if opt.All {
		server, err := t.GetServer()
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
	return StringSliceDiff(names, opt.Exclude), nil
}
