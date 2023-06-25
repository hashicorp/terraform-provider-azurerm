package cluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterUpgradePolicy struct {
	DeltaHealthPolicy             *ClusterUpgradeDeltaHealthPolicy `json:"deltaHealthPolicy,omitempty"`
	ForceRestart                  *bool                            `json:"forceRestart,omitempty"`
	HealthCheckRetryTimeout       string                           `json:"healthCheckRetryTimeout"`
	HealthCheckStableDuration     string                           `json:"healthCheckStableDuration"`
	HealthCheckWaitDuration       string                           `json:"healthCheckWaitDuration"`
	HealthPolicy                  ClusterHealthPolicy              `json:"healthPolicy"`
	UpgradeDomainTimeout          string                           `json:"upgradeDomainTimeout"`
	UpgradeReplicaSetCheckTimeout string                           `json:"upgradeReplicaSetCheckTimeout"`
	UpgradeTimeout                string                           `json:"upgradeTimeout"`
}
