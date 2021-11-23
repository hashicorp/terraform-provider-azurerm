package application

type ApplicationUpgradePolicy struct {
	ApplicationHealthPolicy        *ApplicationHealthPolicy        `json:"applicationHealthPolicy,omitempty"`
	ForceRestart                   *bool                           `json:"forceRestart,omitempty"`
	InstanceCloseDelayDuration     *int64                          `json:"instanceCloseDelayDuration,omitempty"`
	RecreateApplication            *bool                           `json:"recreateApplication,omitempty"`
	RollingUpgradeMonitoringPolicy *RollingUpgradeMonitoringPolicy `json:"rollingUpgradeMonitoringPolicy,omitempty"`
	UpgradeMode                    *RollingUpgradeMode             `json:"upgradeMode,omitempty"`
	UpgradeReplicaSetCheckTimeout  *int64                          `json:"upgradeReplicaSetCheckTimeout,omitempty"`
}
