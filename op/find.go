/* SPDX-License-Identifier: Apache-2.0
*  Copyright 2018 Alex Athanasopoulos
*/
package op

import (
	"fmt"
	"os"
	"strconv"

	"github.com/melato/lxdtool/proc"
)

func (c *Server) GetPidMap() (map[int]string, error) {
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

func (t *Server) FindPid(ProcDir string, pmap map[int]string, pid int) error {
	ps := proc.NewProc(ProcDir)
	p := pid
	var stat *proc.Stat
	for {
		if p == 1 {
			return nil
		}
		s, err := ps.GetStat(p)
		if err != nil {
			return err
		}
		if p == pid {
			stat = s
		}
		name, ok := pmap[p]
		if ok {
			fmt.Println(pid, name, stat.Name)
			return nil
		}
		p = s.Ppid
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("process", proc)
	return nil
}

func (t *Server) FindPids(ProcDir string, args []string) error {
	pmap, err := t.GetPidMap()
	if err != nil {
		return err
	}
	for _, s := range args {
		pid, err := strconv.Atoi(s)
		if err == nil {
			t.FindPid(ProcDir, pmap, pid)
		} else {
			fmt.Println("not an int: " + s)
		}
	}
	return nil
}
