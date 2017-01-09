package emanage

import (
	"fmt"
	"time"

	"github.com/elastifile/emanage-go/pkg/etime"
	"github.com/elastifile/emanage-go/pkg/rest"
	"github.com/elastifile/emanage-go/pkg/size"
)

const sysStatsUri = "api/system_statistics"

type statistics struct {
	conn *rest.Session
}

type Statistic struct {
	Id                   StatEntryId             `json:"id"`
	ReadNumEvents        uint                    `json:"read_num_events"`
	MdReadNumEvents      uint                    `json:"md_read_num_events"`
	MdWriteNumEvents     uint                    `json:"md_write_num_events"`
	WriteNumEvents       uint                    `json:"write_num_events"`
	HostsUp              uint                    `json:"hosts_up"`
	HostsTotal           uint                    `json:"hosts_total"`
	DevicesUp            uint                    `json:"devices_up"`
	DevicesTotal         uint                    `json:"devices_total"`
	ActiveConnections    uint                    `json:"active_connections"`
	TimeStamp            time.Time               `json:"timestamp"`
	DataReductionRatio   float64                 `json:"data_reduction_ratio"`
	Connections          statConnByType          `json:"active_connections_by_type"`
	DataReductions       statDataReductionByType `json:"data_reduction_ratio_by_type"`
	ReadIo               statIo                  `json:"read_io"`
	WriteIo              statIo                  `json:"write_io"`
	Capacity             statCapacity            `json:"capacity"`
	FreeCapacityInBytes  statCapacity            `json:"free"`
	UtilizedCapacity     statCapacity            `json:"utilization"`
	SelfHealingWatermark statCapacity            `json:"self_healing_watermark"`
	TotalMemory          statMemory              `json:"total_memory"`
	FreeMemory           statMemory              `json:"free_memory"`
	ReadMaxLatency       statLatency             `json:"read_max_latency"`
	ReadLatency          statLatency             `json:"read_latency"`
	ReadMinLatency       statLatency             `json:"read_min_latency"`
	MdReadLatency        statLatency             `json:"md_read_latency"`
	MdWriteLatency       statLatency             `json:"md_write_latency"`
	WriteLatency         statLatency             `json:"write_latency"`
	WriteMaxLatency      statLatency             `json:"write_max_latency"`
	WriteMinLatency      statLatency             `json:"write_min_latency"`
	EnodeStatistics      []statEnode             `json:"enode_statistics"`
}

func (s Statistic) FreeCapacity() size.Size {
	return size.Size(s.FreeCapacityInBytes.Bytes) * size.Byte
}

func (s *Statistic) GetUsedPercentage() (float64, error) {
	freeCapacity := s.FreeCapacityInBytes.Bytes
	capacity := s.Capacity.Bytes
	selfHealingWM := s.SelfHealingWatermark.Bytes
	result := 0.0

	Log.Info("GetUsedPrecentage",
		"Free Capacity", freeCapacity,
		"Capacity", capacity,
		"Self Healing WM", selfHealingWM,
	)
	if selfHealingWM > 0 {
		result = 100 * (1.0 - (float64(freeCapacity-capacity+selfHealingWM) / float64(selfHealingWM)))
	}
	return result, nil
}

type statConnByType struct {
	Nfs    uint `json:"NFS"` // FIXME: not according to lower case standard
	Smb    uint `json:"SMB"` // FIXME: not according to lower case standard
	Block  uint `json:"block"`
	Object uint `json:"object"`
}

type statDataReductionByType struct {
	Dedup       float64 `json:"Dedup"`       // FIXME: not according to lower case standard
	Compression float64 `json:"Compression"` // FIXME: not according to lower case standard
}

type statIo struct {
	Bytes int64 `json:"bytes"`
}

type statCapacity struct {
	// In some calculations we can get negative values, so we need this to be signed
	Bytes int64 `json:"bytes"`
}

type statMemory struct {
	Bytes interface{} `json:"bytes"` // FIXME: this was null
}

type statLatency struct {
	Nanos int `json:"nanos"` // FIXME: this is int because read_min_latency was -1 ?????
}

func (s *statistics) GetAll(startTime etime.NilableTime) ([]Statistic, error) {
	param := struct {
		StartTime etime.NilableTime `json:"start_time,omitempty"`
	}{startTime}
	var result []Statistic
	if err := s.conn.Request(rest.MethodGet, sysStatsUri, param, &result); err != nil {
		return nil, err
	}

	return result, nil
}

type statFullRes struct {
	Statistic
	// EnodeStatistics    []statEnode             `json:"enode_statistics"`
	AverageReplication uint `json:"average_replication"`
}

type statEnode struct {
	EnodeId           uint           `json:"enode_id"`
	Status            string         `json:"status"` //TODO: define a type
	ReadNumEvents     uint           `json:"read_num_events"`
	WriteNumEvents    uint           `json:"write_num_events"`
	ActiveConnections uint           `json:"active_connections"`
	Uptime            uint           `json:"uptime"`
	Timestamp         time.Time      `json:"timestamp"`
	Connections       statConnByType `json:"active_connections_by_type"`
	CpuUserTime       statCpu        `json:"cpu_user_time"`
	CpuSystemTime     statCpu        `json:"cpu_system_time"`
	TotalMemory       statMemory     `json:"total_memory"`
	FreeMemory        statMemory     `json:"free_memory"`
	ReadIo            statIo         `json:"read_io"`
	WriteIo           statIo         `json:"write_io"`
	ReadMaxLatency    statLatency    `json:"read_max_latency"`
	ReadLatency       statLatency    `json:"read_latency"`
	ReadMinLatency    statLatency    `json:"read_min_latency"`
	WriteLatency      statLatency    `json:"write_latency"`
	WriteMaxLatency   statLatency    `json:"write_max_latency"`
	WriteMinLatency   statLatency    `json:"write_min_latency"`
	Enode             statEnodeInner `json:"enode"`
}

type statCpu struct {
	Percent float64 `json:"percent"` // TODO: maybe make into a dedicated precntage type ...
}

type statEnodeInner struct { // FIXME: why is it exists? why this nesting is needed?
	Cores           uint        `json:"cores"`
	Memory          interface{} `json:"memory"`           // FIXME: this was null
	Role            string      `json:"role"`             // TODO: should make a type for that
	Status          string      `json:"status"`           // TODO: should make a type for that
	SoftwareVersion interface{} `json:"software_version"` // FIXME: this was null
	Host            statHost    `json:"host"`
}

type statHost struct {
	Id                 uint        `json:"id"`
	Name               string      `json:"name"`
	Vendor             interface{} `json:"vendor"`        // FIXME: this was null
	Model              interface{} `json:"model"`         // FIXME: this was null
	Path               interface{} `json:"path"`          // FIXME: this was null
	Status             string      `json:"status"`        // FIXME: what the string option exists? saw "0" value ?
	PowerState         interface{} `json:"power_state"`   // FIXME: this was null
	Cpus               interface{} `json:"cpus"`          // FIXME: this was null
	MemoryMb           interface{} `json:"memory_mb"`     // FIXME: this was null
	Software           interface{} `json:"software"`      // FIXME: this was null
	Maintenance        interface{} `json:"maintenance"`   // FIXME: this was null
	Nincs              interface{} `json:"nics"`          // FIXME: this was null
	VmManagerId        interface{} `json:"vm_manager_id"` // FIXME: this was null
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	AvailabilityZoneId interface{} `json:"availability_zone_id"` // FIXME: this was null
	DistributionZoneId interface{} `json:"distribution_zone_id"` // FIXME: this was null
}

func (s *statistics) GetFull(entryId StatEntryId) (statFullRes, error) {
	uri := fmt.Sprintf("%s/%d", sysStatsUri, entryId)
	var result statFullRes
	err := s.conn.Request(rest.MethodGet, uri, nil, &result)
	return result, err
}
