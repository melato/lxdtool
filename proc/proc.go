package proc

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func Getppid(pid int) (int, error) {
	stat, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
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
