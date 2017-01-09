package emanage

type EmEvents int

const (
	DiskIOErr                               EmEvents = 1
	TransDiskErrorWithMsg                   EmEvents = 2
	ClientDisconnect                        EmEvents = 3
	ClientConnect                           EmEvents = 4
	MPoolResourceLimitExceeded              EmEvents = 5
	EThreadsResourceLimitExceeded           EmEvents = 6
	RPCHandlerResourceLimitExceeded         EmEvents = 7
	PrstntStoreSvcConcurLimit               EmEvents = 8
	PstoreDirectBucketResourceLimitExceeded EmEvents = 9
	InternalConnectivityFailure             EmEvents = 10
	ControlSystemStarted                    EmEvents = 11
	ClienFailedToMount                      EmEvents = 12
	Panic                                   EmEvents = 13
	LogEvent                                EmEvents = 14
	DiskFormatted                           EmEvents = 15
	SlowDevice                              EmEvents = 16
	ConnectToEcsFailed                      EmEvents = 17
	ClctNewDataDistributionMap              EmEvents = 1003
	TrnasferToMapCompleted                  EmEvents = 1004
	NodeAdded                               EmEvents = 1005
	NodeRemoved                             EmEvents = 1006
	DiskAdded                               EmEvents = 1007
	DiskRemoved                             EmEvents = 1008
	DiskIntegratedIn                        EmEvents = 1009
	DiskIntegratedOut                       EmEvents = 1010
	NicAdded                                EmEvents = 1011
	NicRemoved                              EmEvents = 1012
	ClusterProtocolServiceStopped           EmEvents = 1013
	ClusterMetadataServiceStopped           EmEvents = 1014
	ClusterStorageServiceStopped            EmEvents = 1015
	ClusterProtocolServiceStarted           EmEvents = 1016
	ClusterMetadataServiceStarted           EmEvents = 1017
	ClusterStorageServiceStarted            EmEvents = 1018
	OwnershipMapDeployed                    EmEvents = 1019
	VheadFenced                             EmEvents = 1020
	ClusterLockdown                         EmEvents = 1021
	ClusterRecovered                        EmEvents = 1022
	VheadErrorState                         EmEvents = 1023
	VheadRecovered                          EmEvents = 1024
	DataObjMapDistributionStarted           EmEvents = 1025
	DataObjMapDistributionFinished          EmEvents = 1026
	DataObjMapDistributionFailed            EmEvents = 1027
	OwnershipMapDistributionStarted         EmEvents = 1028
	OwnershipMapDistributionFinished        EmEvents = 1029
	OwnershipMapDistributionFailed          EmEvents = 1030
	VheadErrTryRecover                      EmEvents = 1031
	VheadRecoverInFEOnlyMode                EmEvents = 1032
	SystemInDegradDataRplctLevel            EmEvents = 1033
	SystemDataRplctLvlRstrdToNormal         EmEvents = 1034
	SystemCapacityExceeded                  EmEvents = 1035
	SystemCapacityBack                      EmEvents = 1036
	UnstableAfterRmvTaskFailure             EmEvents = 1037
	OwnershipRecoveryStarted                EmEvents = 1038
	OwnershipRecoveryFinished               EmEvents = 1039
	EcdbMapDistributionStarted              EmEvents = 1040
	EcdbMapDistributionFinished             EmEvents = 1041
	DcHardQuotaExceeded                     EmEvents = 5005
	DcSoftQuotaExceeded                     EmEvents = 5006
	VheadAdded                              EmEvents = 5007
	VheadRemoved                            EmEvents = 5008
	VheadRemovalStarted                     EmEvents = 5009
	DcCreated                               EmEvents = 5010
	DcRemoved                               EmEvents = 5011
	ExportCreated                           EmEvents = 5012
	ExportRemoved                           EmEvents = 5013
	AddedPermissionToExport                 EmEvents = 5014
	RemovedPermissionToExport               EmEvents = 5015
	LinkDown                                EmEvents = 5016
	LinkUp                                  EmEvents = 5017
	SystemCapacityLow                       EmEvents = 5018
	FailedReadHostsVmStatus                 EmEvents = 5019
)

func (e EmEvents) String() string {
	switch e {
	case DiskIOErr:
		return "Disk IO error"
	case TransDiskErrorWithMsg:
		return "Transient disk error in  with message"
	case ClientDisconnect:
		return "Client disconnected"
	case ClientConnect:
		return "Client connected"
	case MPoolResourceLimitExceeded:
		return "MPoolResourceLimitExceeded"
	case EThreadsResourceLimitExceeded:
		return "EThreadsResourceLimitExceeded"
	case RPCHandlerResourceLimitExceeded:
		return "RPCHandlerResourceLimitExceeded"
	case PrstntStoreSvcConcurLimit:
		return "ersistent store service concurrency limit"
	case PstoreDirectBucketResourceLimitExceeded:
		return "PstoreDirectBucketResourceLimitExceeded"
	case InternalConnectivityFailure:
		return "Internal connectivity failure"
	case ControlSystemStarted:
		return "Control system started"
	case ClienFailedToMount:
		return "Client failed to mount on"
	case Panic:
		return "PANIC"
	case LogEvent:
		return "Log Event"
	case DiskFormatted:
		return "DiskFormatted"
	case SlowDevice:
		return "Slow device"
	case ConnectToEcsFailed:
		return "Connection to Elastifile Control Service (ECS) failed"
	case ClctNewDataDistributionMap:
		return "Calculated new data distribution map"
	case TrnasferToMapCompleted:
		return "Transfer to map completed"
	case NodeAdded:
		return "Node Added"
	case NodeRemoved:
		return "Node Removed"
	case DiskAdded:
		return "Disk Added"
	case DiskRemoved:
		return "Disk Removed"
	case DiskIntegratedIn:
		return "Disk Integrated In"
	case DiskIntegratedOut:
		return "Disk Integrated Out"
	case NicAdded:
		return "Nic Added"
	case NicRemoved:
		return "Nic Removed"
	case ClusterProtocolServiceStopped:
		return "ClusterProtocolServiceStopped"
	case ClusterMetadataServiceStopped:
		return "ClusterMetadataServiceStopped"
	case ClusterStorageServiceStopped:
		return "ClusterStorageServiceStopped"
	case ClusterProtocolServiceStarted:
		return "ClusterProtocolServiceStarted"
	case ClusterMetadataServiceStarted:
		return "ClusterMetadataServiceStarted"
	case ClusterStorageServiceStarted:
		return "ClusterStorageServiceStarted"
	case OwnershipMapDeployed:
		return "Ownership Map Deployed"
	case VheadFenced:
		return "vHead  fenced"
	case ClusterLockdown:
		return "Cluster Lockdown"
	case ClusterRecovered:
		return "Cluster Recovered"
	case VheadErrorState:
		return "VheadIn  Error State"
	case VheadRecovered:
		return "Vhead Recovered"
	case DataObjMapDistributionStarted:
		return "Data Object Map Distribution Started"
	case DataObjMapDistributionFinished:
		return "Data Object Map Distribution Finished"
	case DataObjMapDistributionFailed:
		return "Data Object Map Distribution Failed"
	case OwnershipMapDistributionStarted:
		return "Ownership Map Distribution Started"
	case OwnershipMapDistributionFinished:
		return "Ownership Map Distribution Finished"
	case OwnershipMapDistributionFailed:
		return "Ownership Map Distribution Failed"
	case VheadErrTryRecover:
		return "vHead has encountered an error. Trying to recover"
	case VheadRecoverInFEOnlyMode:
		return "vHead recovered as frontend only node"
	case SystemInDegradDataRplctLevel:
		return "System is in degraded data replication level"
	case SystemDataRplctLvlRstrdToNormal:
		return "System data replication level restored to normal"
	case SystemCapacityExceeded:
		return "System capacity exceeded. All writes are disabled, please free some space or expand cluster capacity"
	case SystemCapacityBack:
		return "System capacity is back operational state. Writes are enabled"
	case OwnershipRecoveryStarted:
		return "Ownership Recovery Started"
	case OwnershipRecoveryFinished:
		return "Ownership Recovery Finished"
	case UnstableAfterRmvTaskFailure:
		return "Unstable after remove task failure"
	case DcHardQuotaExceeded:
		return "Data container hard quota limit exceeded"
	case DcSoftQuotaExceeded:
		return "Data container soft quota limit exceeded"
	case VheadAdded:
		return "Vhead addded to the cluster"
	case VheadRemoved:
		return "Vhead removed from the cluster"
	case VheadRemovalStarted:
		return "vHead removal process started"
	case DcCreated:
		return "Data Container created with policy"
	case DcRemoved:
		return "Data Container removed"
	case ExportCreated:
		return "Export Created"
	case ExportRemoved:
		return "Export Removed"
	case AddedPermissionToExport:
		return "Added Permission To Export"
	case RemovedPermissionToExport:
		return "Removed Permission To Export"
	case LinkDown:
		return "Cluster nodes connection lost"
	case LinkUp:
		return "Cluster nodes connection established"
	case SystemCapacityLow:
		return "System capacity is low, utilization  capacity"
	case FailedReadHostsVmStatus:
		return "Failed reading hosts and vm status from vCenter"
	case EcdbMapDistributionStarted:
		return "ECDB Map Distribution Started"
	case EcdbMapDistributionFinished:
		return "ECDB Map Distribution Finished"
	default:
		return "Illegal evetn Id"
	}
}
