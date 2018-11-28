/* Copyright 2018 Alex Athanasopoulos
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*   http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
*/
package op

import (
	"fmt"
	"os"
	"strconv"

	"melato.org/lxdtool/proc"
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
		ps := proc.NewProc(ProcDir)
		p, err = ps.Getppid(p)
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
