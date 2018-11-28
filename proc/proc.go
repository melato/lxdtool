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
package proc

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

type Proc struct {
	Dir string
}

func NewProc(dir string) *Proc {
	var proc Proc
	proc.Dir = dir
	return &proc
}

func (t *Proc) Getppid(pid int) (int, error) {
	stat, err := ioutil.ReadFile(path.Join(t.Dir, strconv.Itoa(pid), "stat"))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	line := string(stat)
	// 12021 (lxd) S 1
	fields := strings.Fields(line)
	if len(fields) >= 4 {
		ppid, err := strconv.Atoi(fields[3])
		if err != nil {
			return 0, err
		}
		return int(ppid), nil
	}
	return 0, nil
}
