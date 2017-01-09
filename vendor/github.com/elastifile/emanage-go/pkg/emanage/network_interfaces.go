package emanage

import (
	"fmt"
	"path"

	"github.com/elastifile/errors"

	"github.com/elastifile/emanage-go/pkg/rest"
)

var (
	netInterfacesUri = "/api/network_interfaces"
)

type netInterfaces struct {
	conn *rest.Session
}

func (n *netInterfaces) GetAll() ([]networkInterface, error) {
	var result []networkInterface
	if err := n.conn.Request(rest.MethodGet, netInterfacesUri, nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

type UpdateNetInterfacesOpts struct {
	Role string
}

// Updates network_interfaces matching given id with requested UpdateNetInterfacesOpts
func (n *netInterfaces) Update(id int, opts *UpdateNetInterfacesOpts) error {
	Log.Debug("Updating network interface", "netInterface id", id, "opts", opts)
	ifaces, err := n.GetAll()
	if err != nil {
		return err
	}

	type netIfaceOpts struct {
		NetworkInterface *networkInterface `json:"network_interface"`
	}

	for _, iface := range ifaces {
		if iface.ID == id {
			iface.Role = opts.Role
			err := n.conn.Request(rest.MethodPut,
				path.Join(netInterfacesUri, fmt.Sprintf("%v", id)),
				&netIfaceOpts{NetworkInterface: &iface},
				nil)

			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.Errorf("Network interface update: didn't find any network interface with id: %v", id)
}

type networkInterface struct {
	Name           string `json:"name"`
	ID             int    `json:"id"`
	Address        string `json:"address"`
	DetectedByNic1 string `json:"detected_by_nic1"`
	DetectedByNic2 string `json:"detected_by_nic2"`
	Dhcp           bool   `json:"dhcp"`
	EnodeID        int    `json:"enode_id"`
	HostID         int    `json:"host_id"`
	Lladdress      string `json:"lladdress"`
	MaxVfunc       int    `json:"max_vfunc"`
	NumVfunc       int    `json:"num_vfunc"`
	Role           string `json:"role"`
	Speed          int    `json:"speed"`
	SriovActive    bool   `json:"sriov_active"`
	SriovCapable   bool   `json:"sriov_capable"`
	SriovEnabled   bool   `json:"sriov_enabled"`
	Status         string `json:"status"`
	Subnet         string `json:"subnet"`
	URL            string `json:"url"`
	VfuncID        int    `json:"vfunc_id"`
	VswitchID      int    `json:"vswitch_id"`
	VswitchName    string `json:"vswitch_name"`
}
