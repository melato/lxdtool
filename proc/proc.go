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
	"regexp"
	"strconv"
)

type Proc struct {
	Dir string
}

type Stat struct {
	Pid   int
	Name  string
	State string
	Ppid  int
}

func NewProc(dir string) *Proc {
	var proc Proc
	proc.Dir = dir
	return &proc
}

/** should be able to parse things like:
  15063 (ht )tpd) S 10283 ...
*/
var statPattern = regexp.MustCompile(`([0-9]+) \((.*)\) ([A-Za-z]) ([0-9]+)`)

func ParseStat(s string) *Stat {
	fields := statPattern.FindStringSubmatch(s)
	if fields == nil {
		return nil
	}
	var p Stat
	var err error
	p.Pid, err = strconv.Atoi(fields[1])
	if err != nil {
		return nil
	}
	p.Name = fields[2]
	p.State = fields[3]
	p.Ppid, err = strconv.Atoi(fields[4])
	if err != nil {
		return nil
	}
	return &p
}

func (t *Proc) GetStat(pid int) (*Stat, error) {
	data, err := ioutil.ReadFile(path.Join(t.Dir, strconv.Itoa(pid), "stat"))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	line := string(data)
	s := ParseStat(line)
	// 12021 (lxd) S 1
	if s != nil {
		return s, nil
	}
	return nil, nil
}
