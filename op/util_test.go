package op

import (
	"testing"
)

func checkHost(t *testing.T, host string) {
	addr := host + ":1111"
	h := HostAddress(addr)
	if host != h {
		t.Errorf(`HostAddress("%s") returns %s`, addr, h)
	}

}
func TestHostAddress(t *testing.T) {
	checkHost(t, "1.2.3.4")
	checkHost(t, "[::1]")
}
