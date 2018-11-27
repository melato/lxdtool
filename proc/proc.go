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
