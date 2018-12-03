/* SPDX-License-Identifier: Apache-2.0
*  Copyright 2018 Alex Athanasopoulos
*/
package op

import (
	"strings"
)

func HostAddress(address string) string {
	pos := strings.LastIndex(address, ":")
	if pos >= 0 {
		if pos >= 2 && address[0] == '[' && address[pos-1] == ']' {
			return address[1 : pos-1]
		} else {
			return address[0:pos]
		}
	} else {
		return address
	}
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
