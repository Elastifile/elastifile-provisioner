package emanage

import (
	"fmt"
	"path"
	"strings"

	"github.com/elastifile/errors"

	"github.com/elastifile/emanage-go/pkg/rest"
)

var (
	hostsUri     = "/api/hosts"
	syncHostsUri = path.Join(hostsUri, "sync")
)

type hosts struct {
	conn *rest.Session
}

type DetectHostOpts struct {
	Vlan              int    `json:"vlan"`
	HostIDs           []int  `json:"host_ids"`
	BroadcastNic      string `json:"broadcast_nic,omitempty"`
	DataNetworkNumber int    `json:"data_network_number,omitempty"`
}

func (h *hosts) Detect(opts *DetectHostOpts) error {
	Log.Debug("Detecting host connectivity...", "hosts ids", opts.HostIDs, "vlan", opts.Vlan)
	if err := h.conn.Request(rest.MethodPost, path.Join(hostsUri, "detect"), opts, nil); err != nil {
		return err
	}
	return nil
}

func (h *hosts) Sync() error {
	if err := h.conn.Request(rest.MethodPost, syncHostsUri, nil, nil); err != nil {
		return err
	}
	return nil
}

func (h *hosts) GetAll(opt *GetAllOpts) ([]Host, error) {
	if opt == nil {
		opt = &GetAllOpts{}
	}

	var result []Host
	if err := h.conn.Request(rest.MethodGet, hostsUri, opt, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (h *hosts) GetHost(id int) (*Host, error) {
	var result Host
	err := h.conn.Request(
		rest.MethodGet,
		path.Join(hostsUri, fmt.Sprintf("%v", id)),
		nil,
		&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type UpdateHostOpts struct {
	User     string
	Password string
}

// Updates host matching given id with requested UpdateHostOpts
func (h *hosts) Update(id int, opts *UpdateHostOpts) error {
	Log.Debug("Updating host", "host id", id, "opts", opts)
	hosts, err := h.GetAll(nil)
	if err != nil {
		return err
	}

	var result Host

	for _, hst := range hosts {
		if hst.ID == id {

			hst.User = opts.User
			hst.Password = opts.Password
			err := h.conn.Request(
				rest.MethodPut,
				path.Join(hostsUri, fmt.Sprintf("%v", hst.ID)),
				&hostOpts{Host: &hst},
				&result)

			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.Errorf("Host update: didn't find any host with id: %v", id)
}

type Bytes struct {
	Bytes int `json:"bytes"`
}

type Device struct {
	CanonicalName string      `json:"canonical_name"`
	Capacity      Bytes       `json:"capacity"`
	CreatedAt     string      `json:"created_at"`
	DevicePath    string      `json:"device_path"`
	EnodeID       int         `json:"enode_id"`
	Format        interface{} `json:"format"`
	Free          interface{} `json:"free"`
	HostID        int         `json:"host_id"`
	ID            int         `json:"id"`
	IsWritable    interface{} `json:"is_writable"`
	Model         string      `json:"model"`
	Name          string      `json:"name"`
	PciID         string      `json:"pci_id"`
	Ssd           interface{} `json:"ssd"`
	Status        string      `json:"status"`
	UpdatedAt     string      `json:"updated_at"`
	Usage         Bytes       `json:"usage"`
	UUID          string      `json:"uuid"`
	Vendor        string      `json:"vendor"`
	VMID          interface{} `json:"vm_id"`
}

type Host struct {
	Cores                  int                `json:"cores"`
	Datastores             []DataStore        `json:"datastores"`
	Devices                []Device           `json:"devices"`
	DevicesCount           int                `json:"devices_count"`
	EnableSriov            interface{}        `json:"enable_sriov"`
	ID                     int                `json:"id"`
	Maintenance            bool               `json:"maintenance"`
	Memory                 int                `json:"memory"`
	Model                  string             `json:"model"`
	Name                   string             `json:"name"`
	NetworkInterfaces      []networkInterface `json:"network_interfaces"`
	NetworkInterfacesCount int                `json:"network_interfaces_count"`
	Networks               []struct {
		CreatedAt string `json:"created_at"`
		HostID    int    `json:"host_id"`
		ID        int    `json:"id"`
		Name      string `json:"name"`
		UpdatedAt string `json:"updated_at"`
		Vlan      int    `json:"vlan"`
		VswitchID int    `json:"vswitch_id"`
	} `json:"networks"`
	Path        string      `json:"path"`
	PowerState  string      `json:"power_state"`
	Role        string      `json:"role"`
	Software    string      `json:"software"`
	Status      string      `json:"status"`
	User        interface{} `json:"user"`
	Password    interface{} `json:"password"`
	Vendor      string      `json:"vendor"`
	VMManagerID int         `json:"vm_manager_id"`
}

func (h *Host) GetDataStoreByPrefix(prefix string) (*DataStore, error) {
	if len(prefix) == 0 {
		return nil, errors.Errorf("DataStore name is empty")
	}

	for _, d := range h.Datastores {
		if strings.HasPrefix(d.Name, prefix) {
			// found = true
			return &d, nil
		}
	}
	return nil, errors.Errorf("Failed to find datastore with prefix '%v'.", prefix)
}

func (h *Host) DevicesIDsByPrefix(prefix ...string) (result []int) {
	for _, dev := range h.Devices {
		for _, prefix := range prefix {
			if strings.HasPrefix(dev.Name, prefix) {
				Log.Debug("Setup device match", "disk name", dev.Name, "prefix", prefix)
				result = append(result, dev.ID)
			}
		}
	}
	return
}

type hostOpts struct {
	Host *Host `json:"host"`
}
