package emanage

import (
	"fmt"
	"path"
	"time"

	"github.com/elastifile/emanage-go/pkg/eurl"
	"github.com/elastifile/emanage-go/pkg/optional"
	"github.com/elastifile/emanage-go/pkg/rest"

	"github.com/pborman/uuid"
)

const dcUri = "api/data_containers"

type dataContainers struct {
	conn *rest.Session
}

type DataContainer struct {
	Id             DcId      `json:"id"`
	Name           string    `json:"name"`
	Uuid           uuid.UUID `json:"uuid"`
	Used           Bytes     `json:"used_capacity"`
	Scope          string    `json:"namespace_scope"`
	Policy         Policy    `json:"policy"`
	PolicyId       uint      `json:"policy_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Url            eurl.URL  `json:"url"`
	Exports        []Export  `json:"exports"`
	SoftQuota      Bytes     `json:"soft_quota"`
	HardQuota      Bytes     `json:"hard_quota"`
	DirPermissions int       `json:"dir_permissions,omitempty"`
}

type DcGetAllOpts struct {
	GetAllOpts

	FilterByVm     optional.Int `json:"filter_by_vm,omitempty"`     // Represents a vm_id
	FilterByPolicy optional.Int `json:"filter_by_policy,omitempty"` // Represents a policy_id
}

func (dcs *dataContainers) GetAll(opt *DcGetAllOpts) (result []DataContainer, err error) {
	if opt == nil {
		opt = &DcGetAllOpts{}
	}

	if err = dcs.conn.Request(rest.MethodGet, dcUri, opt, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (dcs *dataContainers) GetFull(dcId DcId) (result DataContainer, err error) {
	uri := fmt.Sprintf("%s/%d", dcUri, dcId)
	err = dcs.conn.Request(rest.MethodGet, uri, nil, &result)
	return result, err
}

type DcCreateOpts struct {
	SoftQuota      int             `json:"soft_quota"`
	HardQuota      int             `json:"hard_quota"`
	DirPermissions int             `json:"dir_permissions,omitempty"`
	Share          optional.String `json:"share,omitempty"`
	VmIds          []int           `json:"vm_ids,omitempty"`
}

type DcUpdateOpts struct {
	SoftQuota int             `json:"soft_quota"`
	HardQuota int             `json:"hard_quota"`
	Share     optional.String `json:"share,omitempty"`
}

func (dcs *dataContainers) Create(name string, policyId PolicyId, opt *DcCreateOpts) (DataContainer, error) {
	if opt == nil {
		opt = &DcCreateOpts{}
	}

	params := struct {
		Name     string   `json:"name"`
		PolicyId PolicyId `json:"policy_id"`
		DcCreateOpts
	}{name, policyId, *opt}
	var result DataContainer
	err := dcs.conn.Request(rest.MethodPost, dcUri, params, &result)
	return result, err
}

func (dcs *dataContainers) Update(dc *DataContainer, opt *DcUpdateOpts) (DataContainer, error) {
	if opt == nil {
		panic(fmt.Errorf("requireing update opts"))
	}

	params := struct {
		Name string `json:"name"`
		DcUpdateOpts
	}{dc.Name, *opt}
	var result DataContainer
	uri := fmt.Sprintf("%s/%d", dcUri, dc.Id)
	err := dcs.conn.Request(rest.MethodPut, uri, params, &result)
	return result, err
}

func (dcs *dataContainers) Delete(dc *DataContainer) (result DataContainer, err error) {
	uri := path.Join(dcUri, fmt.Sprintf("%v", dc.Id))
	result = DataContainer{}
	if err = dcs.conn.Request(rest.MethodDelete, uri, nil, &result); err != nil {
		return result, err
	}

	return result, nil
}
