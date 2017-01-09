package emanage

import (
	"path"

	"github.com/elastifile/emanage-go/pkg/rest"
)

var (
	vmsUri     = "/api/vms"
	syncVMsUri = path.Join(vmsUri, "sync")
)

type VM struct {
	Name   string `json:"name"`
	Cores  int    `json:"cores"`
	HostID int    `json:"host_id"`

	Disks []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"devices"`

	Networks []struct {
		MAC  string `json:"mac"`
		IP   string `json:"ip,omitempty"`
		Name string `json:"name"`
	} `json:"networks"`
}

type vms struct {
	conn *rest.Session
}

func (vms *vms) Sync() error {
	if err := vms.conn.Request(rest.MethodPost, syncVMsUri, nil, nil); err != nil {
		return err
	}
	return nil
}

func (vms *vms) GetAll() ([]VM, error) {
	var result []VM
	if err := vms.conn.Request(rest.MethodGet, vmsUri, nil, &result); err != nil {
		return result, err
	}
	return result, nil
}
