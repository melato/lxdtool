package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"melato.org/lxdtool/op"
)

/*
	snapserver -r star -c ~/snap/lxd/current/.config/lxc -p 8080
*/

type SnapServer struct {
	Server op.Server
	Addr   string
}

func (t *SnapServer) Error(w http.ResponseWriter, err error) {
	fmt.Println(err)
	http.Error(w, "Internal Error", 404)
	return
}

func (t *SnapServer) findContainerFromAddress(addr string) (string, error) {
	fields := strings.Split(addr, ":")
	ip := fields[0]

	server, err := t.Server.GetServer()
	if err != nil {
		return "", err
	}
	containers, err := server.GetContainers()
	if err != nil {
		return "", err
	}
	for _, container := range containers {
		if container.IsActive() {
			state, _, err := server.GetContainerState(container.Name)
			if err != nil {
				return "", err
			}
			for _, network := range state.Network {
				for _, a := range network.Addresses {
					if ip == a.Address {
						return container.Name, nil
					}
				}
			}
		}
	}
	return "", nil
}

func (t *SnapServer) Id(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Not implemented", 501)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	body := make(map[string]interface{})

	var err error
	body["RemoteAddr"] = r.RemoteAddr
	body["Name"], err = t.findContainerFromAddress(r.RemoteAddr)
	if err != nil {
		t.Error(w, err)
		return
	}
	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		http.Error(w, "Internal server error", 500)
		return
	}
}

func (t *SnapServer) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Not implemented", 501)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	server, err := t.Server.GetServer()
	if err != nil {
		t.Error(w, err)
		return
	}

	name, err := t.findContainerFromAddress(r.RemoteAddr)
	if err != nil {
		t.Error(w, err)
		return
	}
	snapshots, err := server.GetContainerSnapshotNames(name)
	if err != nil {
		t.Error(w, err)
		return
	}
	body := make(map[string]interface{})
	body["snapshots"] = snapshots
	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		http.Error(w, "Internal server error", 500)
		return
	}
}

func (t *SnapServer) Run() error {
	r := mux.NewRouter()
	r.HandleFunc("/1.0/id", func(w http.ResponseWriter, r *http.Request) { t.Id(w, r) })
	r.HandleFunc("/1.0/list", func(w http.ResponseWriter, r *http.Request) { t.List(w, r) })

	fmt.Println("starting http server at:", t.Addr)
	err := http.ListenAndServe(t.Addr, r)
	if err != nil {
		return err
	}

	return nil

}

func main() {
	var server SnapServer
	var cmd = &cobra.Command{
		Use:   "snapserver",
		Short: "Handles remote requests from containers, so they can manage their own snapshots.",
		Run:   func(cmd *cobra.Command, args []string) { server.Run() },
	}
	cmd.PersistentFlags().StringVarP(&server.Server.Socket, "socket", "s", "/var/snap/lxd/common/lxd/unix.socket", "LXD unix socket")
	cmd.PersistentFlags().StringVarP(&server.Server.Remote, "remote", "r", "", "LXD remote")
	cmd.PersistentFlags().StringVarP(&server.Server.ConfigDir, "config", "c", "", "lxc config dir (with client.crt)")
	cmd.PersistentFlags().StringVarP(&server.Addr, "listen", "l", ":8080", "listen address")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
