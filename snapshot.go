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

func (t *SnapClient) list() error {
	url := t.BaseUrl + "/1.0/list"
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
	url := t.BaseUrl + "/1.0/create/" + name
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
	url := t.BaseUrl + "/1.0/delete/" + strings.Join(names, ",")
	client := &http.Client{}
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	_, err = ioutil.ReadAll(r.Body)
	return err
}

func (t *SnapClient) RunE(args []string) error {
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
		Use:   "snapshot",
		Short: "create, list, delete snapshots from a container",
		Run:   func(cmd *cobra.Command, args []string) { client.Run(args) },
	}
	cmd.PersistentFlags().StringVarP(&client.BaseUrl, "url", "u", "", "server url")
	cmd.PersistentFlags().BoolVarP(&client.Delete, "delete", "d", false, "delete snapshots")
	cmd.PersistentFlags().BoolVarP(&client.List, "list", "l", false, "list snapshots")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
