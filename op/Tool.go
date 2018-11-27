package op

import (
	"fmt"
	"os"
	"path"

	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/lxc/config"
	"github.com/lxc/lxd/shared/api"
)

type Tool struct {
	ServerSocket string
	ServerRemote string
	ConfigDir    string
	ProcDir      string
	server       lxd.ContainerServer
	All          bool
	Exclude      []string
}

/*
func (t *Tool) loadConfigFile(file string) (string, error) {
	bytes, err := io.ioutil.ReadFile(path.Join(t.ConfigPath, file))
	if err != nil {
		return nil, err
	}
	return []string(bytes), nil
}
			// Connect to LXD over HTTPS
			var args lxd.ConnectionArgs
			args.TLSClientCert, err = loadConfigFile("client.crt")
			if err != nil {
				return nil, err
			}
			args.TLSClientKey, err = loadConfigFile("client.key")
			if err != nil {
				return nil, err
			}
			t.server, err = lxd.ConnectLXD(t.ServerUrl, t.ServerSocket, nil)
*/

func (t *Tool) GetServer() (lxd.ContainerServer, error) {
	if t.server == nil {
		var err error
		if t.ServerRemote != "" && t.ConfigDir != "" {
			fmt.Println("using remote: ", t.ServerRemote)
			fmt.Println("ConfigDir: ", t.ConfigDir)
			confPath := os.ExpandEnv(path.Join(t.ConfigDir, "config.yml"))
			conf, err := config.LoadConfig(confPath)
			if err != nil {
				return nil, err
			}
			t.server, err = conf.GetContainerServer(t.ServerRemote)
			if err != nil {
				return nil, err
			}
		} else {
			// Connect to LXD over the Unix socket
			t.server, err = lxd.ConnectLXDUnix(t.ServerSocket, nil)
			if err != nil {
				return nil, err
			}
		}
	}
	return t.server, nil
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

func (t *Tool) GetContainerNames(args []string) ([]string, error) {
	var names []string
	if t.All {
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
	return StringSliceDiff(names, t.Exclude), nil
}

func (t *Tool) ListContainers(args []string) error {
	names, err := t.GetContainerNames(args)
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
