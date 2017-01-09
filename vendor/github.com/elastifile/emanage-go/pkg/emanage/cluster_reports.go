package emanage

import "github.com/elastifile/emanage-go/pkg/rest"

const sysReportUri = "api/cluster_reports/recent"

type clusterReports struct {
	conn *rest.Session
}

type ClusterReport struct {
	Id                     int    `json:"id"`
	SystemID               int    `json:"system_id"`
	Timestamp              string `json:"timestamp"`
	RocTransitionTotal     int    `json:"roc_transition_total"`
	RocTransitionDone      int    `json:"roc_transition_done"`
	OwnerShipRecoveryTotal int    `json:"ownership_recovery_total"`
	OwnerShipRecoveryDone  int    `json:"ownership_recovery_done"`
}

func (cr *clusterReports) GetAll() (result []ClusterReport, err error) {
	if err = cr.conn.Request(rest.MethodGet, sysReportUri, nil, &result); err != nil {
		Log.Error("GetAll Error", "err", err)
		return nil, err
	}
	return result, nil
}
