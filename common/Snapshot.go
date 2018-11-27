package common

import (
	"time"
)

type Snapshot struct {
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

type Result struct {
	Snapshots []Snapshot `json:"snapshots"`
}
