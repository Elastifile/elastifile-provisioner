package emanage

import (
	"github.com/elastifile/emanage-go/pkg/optional"
)

type (
	StatEntryId uint
	DcId        uint
	ExportId    uint
	PolicyId    uint
	TenantId    uint
	SystemId    uint
)

////////////////////////////////////////

// Source: elastifile_mom.thrift in emanage/ecs

type HealthState string

const (
	HealthOk        HealthState = "Ok"
	HealthAttn      HealthState = "Attn"
	HealthAlert     HealthState = "Alert"
	HealthRisk      HealthState = "Risk"
	HealthCritical  HealthState = "Critical"
	HealthDown      HealthState = "Down"
	HealthEmergency HealthState = "EMERGENCY"
)

var HealthValues = map[string]HealthState{
	"normal":    HealthOk,
	"minor":     HealthAttn,
	"moderate":  HealthAlert,
	"major":     HealthRisk,
	"critical":  HealthCritical,
	"down":      HealthDown,
	"emergency": HealthEmergency,
}

type DedupLevel int
type CompressionLevel int
type ReplicationLevel int

type GetAllOpts struct {
	Search  optional.String `json:"search,omitempty"`
	Order   optional.String `json:"order,omitempty"`
	Page    optional.Int    `json:"page,omitempty"`
	PerPage optional.Int    `json:"per_page,omitempty"`
}

type RecentOpts struct { // move to infra/emanage/simple_types.go ??
	Since int `json:"since,omitempty"` // since task id X
}
