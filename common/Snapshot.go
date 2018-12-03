/* SPDX-License-Identifier: Apache-2.0
*  Copyright 2018 Alex Athanasopoulos
*/
package common

import (
	"time"
)

type Snapshot struct {
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

type Result struct {
	Error     string     `json:"error"`
	Snapshots []Snapshot `json:"snapshots"`
}
