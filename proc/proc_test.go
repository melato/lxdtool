// SPDX-License-Identifier: Apache-2.0
// Copyright 2018 Alex Athanasopoulos
package proc

import (
	"testing"
)

func TestHostAddress(t *testing.T) {
	p := ParseStat("20128 (ht)tpd) S 16390")
	if p.Pid != 20128 {
		t.Errorf("expected pid 20128 actual: %d", p.Pid)
	}
	if p.Ppid != 16390 {
		t.Errorf("expected ppid 16390 actual: %d", p.Ppid)
	}
	if p.Name != "ht)tpd" {
		t.Errorf("expected name ht)tpd actual: %s", p.Name)
	}
}
