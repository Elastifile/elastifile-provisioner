package emanage

import (
	"fmt"
	"path"
	"time"

	"github.com/elastifile/errors"

	"github.com/elastifile/emanage-go/pkg/eurl"
	"github.com/elastifile/emanage-go/pkg/rest"
)

var Timeout time.Duration = 10 * time.Minute

var AfterShutdown func()
var BeforeStart func()
var BeforeForceReset func()

const (
	systemsUri = "api/systems"
	sysID      = 1
)

type systems struct {
	session *rest.Session
}

type SystemState string

const (
	StateSystemInit    SystemState = "system_init"
	StateConfigured    SystemState = "configured"
	StateMapped        SystemState = "mapped"
	StateInService     SystemState = "in_service"
	StateClosingGates  SystemState = "closing_gates"
	StateUpClosedGates SystemState = "up_closed_gates"
	StateShuttingDown  SystemState = "shutting_down"
	StateDown          SystemState = "down"
	StateLockdown      SystemState = "lockdown"
	StateUnknown       SystemState = "unknown"
)

type SystemDetails struct {
	Name             string      `json:"name"`
	Id               SystemId    `json:"id"`
	Status           SystemState `json:"status,omitempty"`
	ConnectionStatus string      `json:"connection_status,omitempty"`
	Uptime           string      `json:"uptime,omitempty"`
	Version          string      `json:"version,omitempty"`
	ReplicationLevel int         `json:"replication_level,omitempty"`
	ControlAddress   string      `json:"control_address,omitempty"`
	ControlPort      int         `json:"control_port,omitempty"`
	NfsAddress       string      `json:"nfs_address,omitempty"`
	NfsIpRange       int         `json:"nfs_ip_range,omitempty"`
	DataAddress      string      `json:"data_address,omitempty"`
	DataIpRange      int         `json:"data_ip_range,omitempty"`
	DataVlan         int         `json:"data_vlan,omitempty"`
	DataAddress2     string      `json:"data_address2,omitempty"`
	DataIpRange2     int         `json:"data_ip_range2,omitempty"`
	DataVlan2        int         `json:"data_vlan2,omitempty"`
	DataMtu          int         `json:"data_mtu,omitempty"`
	DataMtu2         int         `json:"data_mtu2,omitempty"`
	DeploymentModel  string      `json:"deployment_model,omitempty"`
	ExternalUseDhcp  bool        `json:"external_use_dhcp,omitempty"`
	ExternalAddress  string      `json:"external_address,omitempty"`
	ExternalIpRange  string      `json:"external_ip_range,omitempty"`
	ExternalGateway  string      `json:"external_gateway,omitempty"`
	ExternalNetwork  string      `json:"external_network,omitempty"`
	CreatedAt        time.Time   `json:"created_at,omitempty"`
	UpdatedAt        time.Time   `json:"updated_at,omitempty"`
	Url              *eurl.URL   `json:"url,omitempty"`
	TimeZone         string      `json:"time_zone,omitempty"`
	NTPServers       string      `json:"ntp_servers,omitempty"`
}

type StateError struct {
	Expected SystemState
	Actual   SystemState
}

func (e *StateError) Error() string {
	return fmt.Sprintf("Expected system state '%v', Actual is '%v'", e.Expected, e.Actual)
}

func (ss *systems) GetAll(opt *GetAllOpts) ([]SystemDetails, error) {
	if opt == nil {
		opt = &GetAllOpts{}
	}

	var result []SystemDetails
	return result, ss.session.Request(rest.MethodGet, systemsUri, opt, &result)
}

func (ss *systems) GetById(id SystemId) (*System, *SystemDetails, error) {
	uri := fmt.Sprintf("%s/%d", systemsUri, id)
	var result SystemDetails
	err := ss.session.Request(rest.MethodGet, uri, nil, &result)
	if err != nil {
		return nil, nil, err
	}
	system := System{
		session: ss.session,
		id:      id,
	}
	return &system, &result, nil
}

func (ss *systems) Update(id SystemId, sysInfo *SystemDetails) (*SystemDetails, error) {
	uri := fmt.Sprintf("%s/%d", systemsUri, id)

	var result SystemDetails
	err := ss.session.Request(rest.MethodPut, uri, &sysInfo, &result)
	return &result, err
}

type System struct {
	session *rest.Session
	id      SystemId
}

func (s *System) anyRequest(method rest.HttpMethod, endpoint string, body interface{}, async bool, result interface{}) error {
	parts := []string{systemsUri, fmt.Sprintf("%d", s.id)}
	if endpoint != "" {
		parts = append(parts, endpoint)
	}
	uri := path.Join(parts...)

	if async {
		return s.session.AsyncRequest(method, uri, body)
	} else {
		return s.session.Request(method, uri, body, result)
	}
}

func (s *System) anyRequestWithDetailsResponse(method rest.HttpMethod, endpoint string, body interface{}, async bool) (*SystemDetails, error) {
	var result SystemDetails

	if err := s.anyRequest(method, endpoint, body, async, &result); err != nil {
		return nil, err
	}

	if async {
		Log.Debug("Received new system state", "details", result)
	}

	return &result, nil
}

func (s *System) request(method rest.HttpMethod, endpoint string, body interface{}) (*SystemDetails, error) {
	return s.anyRequestWithDetailsResponse(method, endpoint, body, false)
}

func (s *System) asyncRequest(method rest.HttpMethod, endpoint string, body interface{}) (*SystemDetails, error) {
	return s.anyRequestWithDetailsResponse(method, endpoint, body, true)
}

func (s *System) GetDetails() (*SystemDetails, error) {
	return s.request(rest.MethodGet, "", nil)
}

func (s *System) ForceReset() (*SystemDetails, error) {
	if BeforeForceReset != nil {
		BeforeForceReset()
	}

	params := struct {
		Async     bool `json:"async"`
		SkipTests bool `json:"skip_tests"`
	}{
		Async:     true,
		SkipTests: true,
	}
	return s.asyncRequest(rest.MethodPost, "force_reset", &params)
}

func (s *System) GetHealth() (*Health, error) {
	var health Health
	healthUri := "health"
	uri := path.Join(systemsUri, fmt.Sprintf("%d", s.id), healthUri)

	err := s.session.Request(rest.MethodGet, uri, nil, &health)
	if err != nil {
		return nil, err
	}
	return &health, nil
}

func (s *System) AcceptEULA() error {
	acceptEULAUri := "accept_eula"
	uri := path.Join(systemsUri, fmt.Sprintf("%d", s.id), acceptEULAUri)
	return s.session.Request(rest.MethodPost, uri, nil, nil)
}

type SystemStartOpts struct {
	SkipTests       bool
	MeltdownRecovey bool
}

func (s *System) Start(opts SystemStartOpts) (*SystemDetails, error) {
	if BeforeStart != nil {
		BeforeStart()
	}

	_, err := s.Setup(nil, opts.SkipTests)
	if err != nil {
		return nil, err
	}

	params := struct {
		Async            bool `json:"async"`
		CreateDefaults   bool `json:"create_defaults"`
		MeltdownRecovery bool `json:"meltdown_recovery"`
	}{
		Async:            true,
		CreateDefaults:   false,
		MeltdownRecovery: opts.MeltdownRecovey,
	}
	return s.asyncRequest(rest.MethodPost, "start", &params)
}

func (s *System) Shutdown() (*SystemDetails, error) {
	details, err := s.asyncRequest(rest.MethodPost, "shutdown", nil)
	if err != nil {
		return nil, err
	}

	if AfterShutdown != nil {
		AfterShutdown()
	}

	return details, nil
}

func (s *System) Setup(answers map[string]interface{}, skipTests bool) (*SystemDetails, error) {
	params := struct {
		Async              bool                   `json:"async"`
		SkipTests          bool                   `json:"skip_tests"`
		ForceResignDevices bool                   `json:"force_resign_devices"`
		AnswerFile         map[string]interface{} `json:"answer_file,omitempty"`
	}{
		Async:              true,
		SkipTests:          skipTests,
		ForceResignDevices: true,
		AnswerFile:         answers,
	}
	return s.asyncRequest(rest.MethodPost, "setup", &params)
}

func (s *System) Deploy() (*SystemDetails, error) {
	return s.asyncRequest(rest.MethodPost, "deploy", nil)
}

func (s *System) AnswerFile() (*AnswerFile, error) {
	answerFileUri := "answer_file"
	uri := path.Join(systemsUri, fmt.Sprintf("%d", s.id), answerFileUri)

	var result AnswerFile
	err := s.session.Request(rest.MethodGet, uri, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type AnswerFile struct {
	Nodes []AnswerFileNode `json:"nodes"`
}

func (af *AnswerFile) ByHost(host string) (n *AnswerFileNode, err error) {
	for _, node := range af.Nodes {
		if node.Name == host {
			nd := node
			n = &nd
		}
	}
	if n == nil {
		return nil, errors.Errorf("Failed mathcing node by name: %s", host)
	}
	return
}

// answerFile body is different than other rest api answers
// in a way that it doesn't have the security prefix.
func (af *AnswerFile) SkipSecurityPrefix() {}

type AnswerFileService struct {
	ID struct {
		Type string `json:"type"`
	} `json:"id"`
	Cores    []int `json:"cores"`
	DeviceID struct {
		UUID string `json:"uuid"`
		Type string `json:"type"`
	} `json:"device_id,omitempty"`
}

type AnswerFileServices []AnswerFileService

func (a AnswerFileServices) ByDeviceID(id ...string) (dstores AnswerFileServices) {
	for _, dstore := range a {
		for _, devID := range id {
			if dstore.DeviceID.UUID == devID {
				dstores = append(dstores, dstore)
			}
		}
	}
	return dstores
}

type AnswerFileNode struct {
	Name     string              `json:"name"`
	Services []AnswerFileService `json:"services"`
}

const ServiceDStore = "DATASTORE_SERVICE_INSTANCE"

func (n AnswerFileNode) DStoreServices() (dstores AnswerFileServices) {
	for _, service := range n.Services {
		if service.ID.Type == ServiceDStore {
			dstores = append(dstores, service)
		}
	}
	return
}

type ReportListElement struct {
	ReportID    string    `json:"report_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IPs         []string  `json:"ips"`
	Time        time.Time `json:"time"`
}

func (s *System) ListReports() ([]ReportListElement, error) {
	var result []ReportListElement
	err := s.anyRequest(rest.MethodGet, "list_reports", nil, false, &result)

	return result, err

}

type CreatedReportDetails struct {
	ID           SystemId  `json:"id"`
	UUID         string    `json:"uuid"`
	LastError    string    `json:"last_error,omitempty"`
	Priority     int       `json:"priority"`
	Attempts     int       `json:"attempts"`
	Queue        string    `json:"queue,omitempty"`
	Name         string    `json:"name"`
	CurrentStep  string    `json:"current_step"`
	StepProgress int       `json:"step_progress"`
	StepTotal    int       `json:"step_total"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Host         string    `json:"host,omitempty"`
	TaskType     string    `json:"task_type"`
	URL          string    `json:"url"`
}

type ReportType string

const (
	ReportTypeFull    ReportType = "full"
	ReportTypeMinimal ReportType = "minimal"
)

func getUUIDsMap(reportList []ReportListElement) map[string]bool {
	m := make(map[string]bool)
	for _, rep := range reportList {
		m[rep.ReportID] = true
	}

	return m
}

func findNewUUIDs(reportListBefore []ReportListElement, reportListAfter []ReportListElement) []string {
	var uuids []string
	beforeUUIDMap := getUUIDsMap(reportListBefore)
	for _, rep := range reportListAfter {
		if _, exists := beforeUUIDMap[rep.ReportID]; !exists {
			uuids = append(uuids, rep.ReportID)
		}
	}

	return uuids
}

func (s *System) CreateReportForNodes(reportType ReportType, ipList []string) ([]string, CreatedReportDetails, error) {
	var result CreatedReportDetails

	reportListBefore, e := s.ListReports()
	if e != nil {
		return nil, result, e
	}

	params := struct {
		Async      bool       `json:"async"`
		ReportType ReportType `json:"report_type"`
		IPList     []string   `json:"ip_list,omitempty"`
	}{
		Async:      true,
		ReportType: reportType,
		IPList:     ipList,
	}
	if err := s.anyRequest(rest.MethodPost, "create_report", &params, true, &result); err != nil {
		return nil, result, err
	}

	reportListAfter, err := s.ListReports()
	if err != nil {
		return nil, result, err
	}

	return findNewUUIDs(reportListBefore, reportListAfter), result, nil
}

func (s *System) CreateReportForAllNodes(reportType ReportType) ([]string, CreatedReportDetails, error) {
	return s.CreateReportForNodes(reportType, nil)
}

func (s *System) DeleteReport(uuid string, ipList []string) (*SystemDetails, error) {
	params := struct {
		ReportID string   `json:"report_id"`
		IPList   []string `json:"ip_list,omitempty"`
	}{
		ReportID: uuid,
		IPList:   ipList,
	}
	return s.request(rest.MethodPost, "delete_report", &params)
}

func (s *System) DeleteReportOnAllNodes(uuid string) (*SystemDetails, error) {
	return s.DeleteReport(uuid, nil)
}

type PreparedReportDetails struct {
	Path string `json:"path,omitempty"`
}

func (s *System) PrepareReport(uuid string, ipList []string) (*PreparedReportDetails, error) {
	var result PreparedReportDetails
	params := struct {
		ReportID string   `json:"report_id"`
		IPList   []string `json:"ip_list,omitempty"`
		PathOnly bool     `json:"path_only,omitempty"`
	}{
		ReportID: uuid,
		IPList:   ipList,
		// TODO: if PathOnly == false The report is sent as a type octet-stream which we do not support
		PathOnly: true,
	}
	if err := s.anyRequest(rest.MethodGet, "download_report", &params, false, &result); err != nil {
		return &result, err
	}

	return &result, nil
}

func (s *System) PrepareReportFromAllNodes(uuid string) (*PreparedReportDetails, error) {
	return s.PrepareReport(uuid, nil)
}

type Capacity struct {
	RawUsage           Bytes  `json:"raw_usage"`
	RawCapacity        Bytes  `json:"raw_capacity"`
	EffectiveUsage     Bytes  `json:"effective_usage"`
	EffectiveCapacity  Bytes  `json:"effective_capacity"`
	DataReductionRatio string `json:"data_reduction_ratio"`
	DedupRatio         string `json:"dedup_ratio"`
	CompressionRatio   string `json:"compression_ratio"`
	TopDataContainers  []struct {
		ID             int    `json:"id"`
		Name           string `json:"name"`
		UUID           string `json:"uuid"`
		UsedCapacity   Bytes  `json:"used_capacity"`
		NamespaceScope string `json:"namespace_scope"`
		DataType       string `json:"data_type"`
		Policy         struct {
			ID          int       `json:"id"`
			Name        string    `json:"name"`
			Dedup       int       `json:"dedup"`
			Compression int       `json:"compression"`
			Replication int       `json:"replication"`
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
			SoftQuota   Bytes     `json:"soft_quota"`
			HardQuota   Bytes     `json:"hard_quota"`
			IsTemplate  bool      `json:"is_template"`
			IsDefault   bool      `json:"is_default"`
		} `json:"policy"`
		PolicyID       int       `json:"policy_id"`
		Dedup          int       `json:"dedup"`
		Compression    int       `json:"compression"`
		SoftQuota      Bytes     `json:"soft_quota"`
		HardQuota      Bytes     `json:"hard_quota"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		ExportsCount   int       `json:"exports_count"`
		DirPermissions int       `json:"dir_permissions"`
		DirUID         int       `json:"dir_uid"`
		DirGid         int       `json:"dir_gid"`
	} `json:"top_data_containers"`
}

func (s *System) Capacity() (*Capacity, error) {
	var result Capacity
	err := s.anyRequest(rest.MethodGet, "capacity", nil, false, &result)
	return &result, err
}
