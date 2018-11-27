package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"melato.org/lxdtool/common"
)

type SnapClient struct {
	BaseUrl string
	Delete  bool
	List    bool
}

type ClientConfig struct {
	BaseUrl string `json:"url"`
}

func (t *SnapClient) list() error {
	url := t.BaseUrl + common.LIST
	client := &http.Client{}
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	var result common.Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	for _, s := range result.Snapshots {
		fmt.Println(s.Date.Format("2006-01-02 15:04:05"), s.Name)
	}
	return nil
}

func (t *SnapClient) create(name string) error {
	url := t.BaseUrl + common.CREATE + "/" + name
	client := &http.Client{}
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	_, err = ioutil.ReadAll(r.Body)
	return err
}

func (t *SnapClient) delete(names []string) error {
	url := t.BaseUrl + common.DELETE + "/" + strings.Join(names, ",")
	client := &http.Client{}
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	_, err = ioutil.ReadAll(r.Body)
	return err
}

func ReadConfig() (*ClientConfig, error) {
	var config ClientConfig
	executable, err := os.Executable()
	if err != nil {
		return &config, err
	}
	configFile := executable + ".json"
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return &config, err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return &config, err
	}
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
