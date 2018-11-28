package op

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lxc/lxd/shared/api"
	"melato.org/lxdtool/common"
)

/*
	snapserver -r star -c ~/snap/lxd/current/.config/lxc -p 8080
*/

type SnapshotServer struct {
	Server *Server
	Addr   string
}

/** Find the container name from its address */
func (t *SnapshotServer) findContainerFromAddress(addr string) (string, error) {
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

type HandlerMethod func(*SnapshotServer, http.ResponseWriter, *http.Request) (
	map[string]interface{}, error)

type HandlerFunction func(http.ResponseWriter, *http.Request)

func (t *SnapshotServer) handler(method HandlerMethod) HandlerFunction {
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

func (t *SnapshotServer) Id(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	body := make(map[string]interface{})
	var err error
	body["RemoteAddr"] = r.RemoteAddr
	body["Name"], err = t.findContainerFromAddress(r.RemoteAddr)
	return body, err
}

func (t *SnapshotServer) List(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	server, err := t.Server.GetServer()
	if err != nil {
		return nil, err
	}

	container, err := t.findContainerFromAddress(r.RemoteAddr)
	if err != nil {
		return nil, err
	}
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
	body := make(map[string]interface{})
	body["snapshots"] = list
	return body, nil
}

func (t *SnapshotServer) Create(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
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

func (t *SnapshotServer) Delete(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
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
