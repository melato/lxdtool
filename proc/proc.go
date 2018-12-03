// SPDX-License-Identifier: Apache-2.0
// Copyright 2018 Alex Athanasopoulos
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
