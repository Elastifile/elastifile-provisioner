package emanage

import (
	"fmt"
	"time"

	"github.com/elastifile/emanage-go/pkg/eurl"
	"github.com/elastifile/emanage-go/pkg/optional"
	"github.com/elastifile/emanage-go/pkg/rest"
)

const policiesUri = "api/policies"

type policies struct {
	conn *rest.Session
}

type Policy struct {
	Id          PolicyId         `json:"id"`
	Name        string           `json:"name"`
	Dedup       DedupLevel       `json:"dedup"`
	Compression CompressionLevel `json:"compression"`
	Replication ReplicationLevel `json:"replication"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Url         eurl.URL         `json:"url"`
}

func (p *policies) GetAll(opt *GetAllOpts) ([]Policy, error) {
	if opt == nil {
		opt = &GetAllOpts{}
	}

	var result []Policy
	if err := p.conn.Request(rest.MethodGet, policiesUri, opt, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *policies) GetFull(policyId PolicyId) (Policy, error) {
	uri := fmt.Sprintf("%s/%d", policiesUri, policyId)
	var result Policy
	err := p.conn.Request(rest.MethodGet, uri, nil, &result)
	return result, err
}

type PolicyCreateOpts struct {
	Dedup       *DedupLevel       `json:"dedup,omitempty"`
	Compression *CompressionLevel `json:"compression,omitempty"`
	Replication *ReplicationLevel `json:"replication,omitempty"`
	TenantId    optional.Int      `json:"tenant_id,omitempty"`
}

func (p *policies) Create(name string, opt *PolicyCreateOpts) (Policy, error) {
	if opt == nil {
		opt = &PolicyCreateOpts{}
	}

	params := struct {
		Name string `json:"name"`
		PolicyCreateOpts
	}{name, *opt}

	var result Policy
	err := p.conn.Request(rest.MethodPost, policiesUri, params, &result)
	return result, err
}
