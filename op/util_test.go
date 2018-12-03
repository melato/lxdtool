/* SPDX-License-Identifier: Apache-2.0
*  Copyright 2018 Alex Athanasopoulos
*/
package op

import (
	"testing"
)

func verifyAddress(t *testing.T, address string, expected string) {
	h := HostAddress(address)
	if expected != h {
		t.Errorf(`HostAddress("%s") returns %s`, address, h)
	}

}
func TestHostAddress(t *testing.T) {
	verifyAddress(t, "1.2.3.4:8080", "1.2.3.4")
	verifyAddress(t, "[::1]:1111", "::1")
}
