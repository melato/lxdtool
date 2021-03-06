// SPDX-License-Identifier: Apache-2.0
// Copyright 2018 Alex Athanasopoulos
package op

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lxc/lxd/shared/api"
	"github.com/melato/lxdtool/common"
)

/*
	snapserver -r star -c ~/snap/lxd/current/.config/lxc -p 8080
*/

type SnapshotServer struct {
	Server  *Server
	Addr    string
	Profile string
}

func (t *SnapshotServer) HasPermission(container *api.Container) bool {
	if t.Profile == "" {
		return true
	}
	for _, p := range container.Profiles {
		if p == t.Profile {
			return true
		}
	}
	return false
}

/** Find the container name from its address */
func (t *SnapshotServer) findContainerFromIP(ip string) (*api.Container, error) {
	server, err := t.Server.GetServer()
	if err != nil {
		return nil, err
	}
	containers, err := server.GetContainers()
	if err != nil {
		return nil, err
	}
	for _, container := range containers {
		if container.IsActive() {
			state, _, err := server.GetContainerState(container.Name)
			if err != nil {
				return nil, err
			}
			for _, network := range state.Network {
				for _, a := range network.Addresses {
					if a.Scope != "local" && ip == a.Address {
						return &container, nil
					}
				}
			}
		}
	}
	return nil, nil
}

func (t *SnapshotServer) Error(w http.ResponseWriter, err error) {
	fmt.Println(err)
	http.Error(w, "Internal Error", 404)
	return
}

func (t *SnapshotServer) start(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return errors.New("Not implemented")
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	return nil
}

type HandlerMethod func(*SnapshotServer, string, http.ResponseWriter, *http.Request) (
	map[string]interface{}, error)

type HandlerFunction func(http.ResponseWriter, *http.Request)

func (t *SnapshotServer) handler(method HandlerMethod) HandlerFunction {
	return func(w http.ResponseWriter, r *http.Request) {
		err := t.start(w, r)
		ip := HostAddress(r.RemoteAddr)
		container, err := t.findContainerFromIP(ip)
		body := make(map[string]interface{})
		if err == nil {
			if container != nil {
				if t.HasPermission(container) {
					body, err = method(t, container.Name, w, r)
				} else {
					body["error"] = "not allowed"
					fmt.Println(ip, "denied")
				}
			} else {
				body["error"] = "not container"
				fmt.Println(ip, "not container")
			}
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

func (t *SnapshotServer) Id(container string, w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	body := make(map[string]interface{})
	body["RemoteAddr"] = r.RemoteAddr
	body["Name"] = container
	return body, nil
}

func (t *SnapshotServer) List(container string, w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	server, err := t.Server.GetServer()
	if err != nil {
		return nil, err
	}
	body := make(map[string]interface{})
	fmt.Println(container, "list")
	snapshots, err := server.GetContainerSnapshots(container)
	if err != nil {
		return nil, err
	}
	var list []common.Snapshot
	for _, snapshot := range snapshots {
		s := common.Snapshot{snapshot.Name, snapshot.CreationDate}
		list = append(list, s)
	}
	body["snapshots"] = list
	return body, nil
}

func (t *SnapshotServer) Create(container string, w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	vars := mux.Vars(r)
	snapshot := vars["snapshot"]

	var err error
	server, err := t.Server.GetServer()
	if err != nil {
		return nil, err
	}

	fmt.Println(container, "create", snapshot)
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

func (t *SnapshotServer) Delete(container string, w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	vars := mux.Vars(r)
	snapshots := strings.Split(vars["snapshots"], ",")
	body := make(map[string]interface{})

	var err error
	server, err := t.Server.GetServer()
	if err != nil {
		return nil, err
	}

	for _, snapshot := range snapshots {
		if snapshot != "" {
			fmt.Println(container, "delete", snapshot)
			err = wait(server.DeleteContainerSnapshot(container, snapshot))
			if err != nil && "not found" != err.Error() {
				return nil, err
			}
		}
	}
	return body, nil
}

func (t *SnapshotServer) Run() error {
	r := mux.NewRouter()
	r.HandleFunc(common.ID, t.handler((*SnapshotServer).Id))
	r.HandleFunc(common.LIST, t.handler((*SnapshotServer).List))
	r.HandleFunc(common.CREATE+"/{snapshot}", t.handler((*SnapshotServer).Create))
	r.HandleFunc(common.DELETE+"/{snapshots}", t.handler((*SnapshotServer).Delete))

	fmt.Println("starting http server at:", t.Addr)
	err := http.ListenAndServe(t.Addr, r)
	if err != nil {
		return err
	}

	return nil

}
