package emanage

import "time"

type Health struct {
	Status  string `json:"status"`
	Details struct {
		ClusterReport struct {
			ID                     int       `json:"id"`
			SystemID               int       `json:"system_id"`
			Timestamp              time.Time `json:"timestamp"`
			RocTransitionTotal     int       `json:"roc_transition_total"`
			RocTransitionDone      int       `json:"roc_transition_done"`
			RocTransitionProgress  int       `json:"roc_transition_progress"`
			OrcTransitionTotal     int       `json:"orc_transition_total"`
			OrcTransitionDone      int       `json:"orc_transition_done"`
			OrcTransitionProgress  int       `json:"orc_transition_progress"`
			EcdbTransitionTotal    int       `json:"ecdb_transition_total"`
			EcdbTransitionDone     int       `json:"ecdb_transition_done"`
			EcdbTransitionProgress int       `json:"ecdb_transition_progress"`
			// RocTotal               interface{} `json:"roc_total"`
			// RocDegraded            interface{} `json:"roc_degraded"`
			// RocRequiredSize        interface{} `json:"roc_required_size"`
			// RocActualSize          interface{} `json:"roc_actual_size"`
			// OrcTotal               interface{} `json:"orc_total"`
			// OrcDegraded            interface{} `json:"orc_degraded"`
			// OrcRequiredSize        interface{} `json:"orc_required_size"`
			// OrcActualSize          interface{} `json:"orc_actual_size"`
			// EcdbTotal              interface{} `json:"ecdb_total"`
			// EcdbDegraded           interface{} `json:"ecdb_degraded"`
			// EcdbRequiredSize       interface{} `json:"ecdb_required_size"`
			// EcdbActualSize         interface{} `json:"ecdb_actual_size"`
			CreatedAt              time.Time `json:"created_at"`
			UpdatedAt              time.Time `json:"updated_at"`
			VersionID              int       `json:"version_id"`
			OwnershipRecoveryTotal int       `json:"ownership_recovery_total"`
			OwnershipRecoveryDone  int       `json:"ownership_recovery_done"`
			RocsWithLastCopy       bool      `json:"rocs_with_last_copy"`
			RocTransitionID        int       `json:"roc_transition_id"`
			OrcTransitionID        int       `json:"orc_transition_id"`
			OwnershipRecoveryID    int       `json:"ownership_recovery_id"`
		} `json:"cluster_report"`
		IsRebuild bool `jhealson:"is_rebuild"`
	} `json:"details"`
}

// {
//   "status": "normal",
//   "details": {
//     "cluster_report": {
//       "id": 236,
//       "system_id": 1,
//       "timestamp": "2016-09-15T05:19:09.000Z",
//       "roc_transition_total": 0,
//       "roc_transition_done": 0,
//       "roc_transition_progress": 0,
//       "orc_transition_total": 0,
//       "orc_transition_done": 0,
//       "orc_transition_progress": 0,
//       "ecdb_transition_total": 0,
//       "ecdb_transition_done": 0,
//       "ecdb_transition_progress": 0,
//       "roc_total": null,
//       "roc_degraded": null,
//       "roc_required_size": null,
//       "roc_actual_size": null,
//       "orc_total": null,
//       "orc_degraded": null,
//       "orc_required_size": null,
//       "orc_actual_size": null,
//       "ecdb_total": null,
//       "ecdb_degraded": null,
//       "ecdb_required_size": null,
//       "ecdb_actual_size": null,
//       "created_at": "2016-09-15T05:19:24.000Z",
//       "updated_at": "2016-09-15T05:19:24.000Z",
//       "version_id": 1,
//       "ownership_recovery_total": 0,
//       "ownership_recovery_done": 0,
//       "rocs_with_last_copy": false,
//       "roc_transition_id": 0,
//       "orc_transition_id": 0,
//       "ownership_recovery_id": 0
//     },
//     "is_rebuild": false
//   }
// }
