/* SPDX-License-Identifier: Apache-2.0
*  Copyright 2018 Alex Athanasopoulos
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/melato/lxdtool/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SnapClient struct {
	BaseUrl string
	Delete  bool
	List    bool
}

type ClientConfig struct {
	BaseUrl string
}

func (t *SnapClient) CallUrl(url string) (*common.Result, error) {
	client := &http.Client{}
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var result common.Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if result.Error != "" {
		return nil, errors.New(result.Error)
	}
	return &result, nil
}

func (t *SnapClient) list() error {
	result, err := t.CallUrl(t.BaseUrl + common.LIST)
	if err != nil {
		return err
	}
	for _, s := range result.Snapshots {
		fmt.Println(s.Date.Format("2006-01-02 15:04:05"), s.Name)
	}
	return nil
}

func (t *SnapClient) create(name string) error {
	_, err := t.CallUrl(t.BaseUrl + common.CREATE + "/" + name)
	return err
}

func (t *SnapClient) delete(names []string) error {
	url := t.BaseUrl + common.DELETE + "/" + strings.Join(names, ",")
	_, err := t.CallUrl(url)
	return err
}

func ReadConfig() (*ClientConfig, error) {
	executable, err := os.Executable()
	if err != nil {
		return nil, err
	}
	configFile := executable + ".yml"

	viper.SetConfigFile(configFile)
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var config ClientConfig
	config.BaseUrl = viper.GetString("url")
	return &config, nil
}

func (t *SnapClient) RunE(args []string) error {
	if t.BaseUrl == "" {
		config, err := ReadConfig()
		if err == nil {
			t.BaseUrl = config.BaseUrl
		}
	}
	if t.BaseUrl == "" {
		return errors.New("missing base url")
	}
	if t.List {
		return t.list()
	} else if t.Delete {
		return t.delete(args)
	} else {
		if len(args) == 1 {
			return t.create(args[0])
		} else {
			return errors.New("must specify exactly one snapshot name")
		}
	}
}

func (t *SnapClient) Run(args []string) {
	err := t.RunE(args)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	var client SnapClient
	var cmd = &cobra.Command{
		Use: "snapshot [flags] [snapshot-name]...",
		Example: `snapshot s1
snapshot -l
snapshot -d s1 s2`,
		Short: "create, list, delete snapshots from inside a container",
		Long: `This program can be run inside a container to manage its own snapshots.
It requires a corresponding snapshot server that relays the requests
to the LXD server.

If run without the -l or -d flags, it creates a snapshot,
deleting any previous snapshot with the same name.
`,
		Run: func(cmd *cobra.Command, args []string) { client.Run(args) },
	}
	cmd.PersistentFlags().StringVarP(&client.BaseUrl, "url", "u", "", "server url")
	cmd.PersistentFlags().BoolVarP(&client.Delete, "delete", "d", false, "delete snapshots")
	cmd.PersistentFlags().BoolVarP(&client.List, "list", "l", false, "list snapshots")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
