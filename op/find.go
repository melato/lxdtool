package op

import (
	"fmt"
	"os"
	"strconv"

	"melato.org/lxdtool/proc"
)

func (c *Tool) GetPidMap() (map[int]string, error) {
	pmap := make(map[int]string)
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
			state, _, err := server.GetContainerState(container.Name)
			if err != nil {
				return nil, err
			}
			pmap[int(state.Pid)] = container.Name
		}
	}
	return pmap, nil
}

func (t *Tool) FindPid(pmap map[int]string, pid int) error {
	p := pid
	for {
		if p == 1 {
			return nil
		}
		name, ok := pmap[p]
		if ok {
			fmt.Println(pid, name)
			return nil
		}
		var err error
		p, err = proc.Getppid(p)
		if err != nil {
			return err
		}
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("process", proc)
	return nil
}

func (t *Tool) FindPids(args []string) error {
	pmap, err := t.GetPidMap()
	if err != nil {
		return err
	}
	for _, s := range args {
		pid, err := strconv.Atoi(s)
		if err == nil {
			t.FindPid(pmap, pid)
		} else {
			fmt.Println("not an int: " + s)
		}
	}
	return nil
}
