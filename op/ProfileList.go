// SPDX-License-Identifier: Apache-2.0
// Copyright 2018 Alex Athanasopoulos
package op

import (
	"fmt"
)

type Profile struct {
	Tool *Tool
}

func (p *Profile) List() error {
	server, err := p.Tool.GetServer()
	if err != nil {
		return err
	}
	names, err := server.GetProfileNames()
	if err != nil {
		return err
	}
	for _, name := range names {
		fmt.Println(name)
	}
	return nil
}
