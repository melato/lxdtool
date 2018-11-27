package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
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

type Snapshot struct {
	Name string
	Date time.Time
}

/** Find the container name from its address */
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

func (t *SnapServer) Error(w http.ResponseWriter, err error) {
	fmt.Println(err)
	http.Error(w, "Internal Error", 404)
	return
}

func (t *SnapServer) start(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return errors.New("Not implemented")
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	return nil
}

type HandlerMethod func(*SnapServer, http.ResponseWriter, *http.Request) (
	map[string]interface{}, error)

type HandlerFunction func(http.ResponseWriter, *http.Request)

func (t *SnapServer) handler(method HandlerMethod) HandlerFunction {
	return func(w http.ResponseWriter, r *http.Request) {
		err := t.start(w, r)
		var body map[string]interface{}
		if err == nil {
			body, err = method(t, w, r)
		}

		if err == nil {
			err = json.NewEncoder(w).Encode(body)
			if err != nil {
				err = errors.New("Internal server error")
			}
		}
		if err != nil {
			t.Error(w, err)
		}
	}
}

func (t *SnapServer) Id(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	body := make(map[string]interface{})
	var err error
	body["RemoteAddr"] = r.RemoteAddr
	body["Name"], err = t.findContainerFromAddress(r.RemoteAddr)
	return body, err
}

func (t *SnapServer) List(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	server, err := t.Server.GetServer()
	if err != nil {
		return nil, err
	}

	container, err := t.findContainerFromAddress(r.RemoteAddr)
	if err != nil {
		return nil, err
	}
	fmt.Println("list", container)
	snapshots, err := server.GetContainerSnapshots(container)
	if err != nil {
		return nil, err
	}
	var list []Snapshot
	for _, snapshot := range snapshots {
		s := Snapshot{snapshot.Name, snapshot.CreationDate}
		list = append(list, s)
	}
	body := make(map[string]interface{})
	body["snapshots"] = list
	return body, nil
}

func wait(op lxd.Operation, err error) error {
	if err == nil {
		return op.Wait()
	}
	return err
}

func (t *SnapServer) Create(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	vars := mux.Vars(r)
	snapshot := vars["snapshot"]

	var err error
	server, err := t.Server.GetServer()
	if err != nil {
		return nil, err
	}

	container, err := t.findContainerFromAddress(r.RemoteAddr)
	if err != nil {
		return nil, err
	}
	fmt.Println("create", container, snapshot)
	err = wait(server.DeleteContainerSnapshot(container, snapshot))
	if err != nil && "not found" != err.Error() {
		return nil, err
	}
	post := api.ContainerSnapshotsPost{
		Name: snapshot,
	}

	err = wait(server.CreateContainerSnapshot(container, post))
	if err != nil {
		return nil, err
	}
	body := make(map[string]interface{})
	body["snapshot"] = snapshot
	return body, nil
}

func (t *SnapServer) Delete(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	vars := mux.Vars(r)
	snapshots := strings.Split(vars["snapshots"], ",")
	body := make(map[string]interface{})
	body["snapshots"] = snapshots

	var err error
	server, err := t.Server.GetServer()
	if err != nil {
		return nil, err
	}

	container, err := t.findContainerFromAddress(r.RemoteAddr)
	if err != nil {
		return nil, err
	}
	for _, snapshot := range snapshots {
		if snapshot != "" {
			fmt.Println("delete", container, snapshot)
			err = wait(server.DeleteContainerSnapshot(container, snapshot))
			if err != nil && "not found" != err.Error() {
				return nil, err
			}
		}
	}
	return body, nil
}

func (t *SnapServer) Run() error {
	r := mux.NewRouter()
	r.HandleFunc("/1.0/id", t.handler((*SnapServer).Id))
	r.HandleFunc("/1.0/list", t.handler((*SnapServer).List))
	r.HandleFunc("/1.0/create/{snapshot}", t.handler((*SnapServer).Create))
	r.HandleFunc("/1.0/delete/{snapshots}", t.handler((*SnapServer).Delete))

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
