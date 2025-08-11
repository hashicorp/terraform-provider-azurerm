package managedcluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterUpgradePolicy struct {
	DeltaHealthPolicy             *ClusterUpgradeDeltaHealthPolicy `json:"deltaHealthPolicy,omitempty"`
	ForceRestart                  *bool                            `json:"forceRestart,omitempty"`
	HealthPolicy                  *ClusterHealthPolicy             `json:"healthPolicy,omitempty"`
	MonitoringPolicy              *ClusterMonitoringPolicy         `json:"monitoringPolicy,omitempty"`
	UpgradeReplicaSetCheckTimeout *string                          `json:"upgradeReplicaSetCheckTimeout,omitempty"`
}
