package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterPropertiesAutoScalerProfile struct {
	BalanceSimilarNodeGroups          *string   `json:"balance-similar-node-groups,omitempty"`
	DaemonsetEvictionForEmptyNodes    *bool     `json:"daemonset-eviction-for-empty-nodes,omitempty"`
	DaemonsetEvictionForOccupiedNodes *bool     `json:"daemonset-eviction-for-occupied-nodes,omitempty"`
	Expander                          *Expander `json:"expander,omitempty"`
	IgnoreDaemonsetsUtilization       *bool     `json:"ignore-daemonsets-utilization,omitempty"`
	MaxEmptyBulkDelete                *string   `json:"max-empty-bulk-delete,omitempty"`
	MaxGracefulTerminationSec         *string   `json:"max-graceful-termination-sec,omitempty"`
	MaxNodeProvisionTime              *string   `json:"max-node-provision-time,omitempty"`
	MaxTotalUnreadyPercentage         *string   `json:"max-total-unready-percentage,omitempty"`
	NewPodScaleUpDelay                *string   `json:"new-pod-scale-up-delay,omitempty"`
	OkTotalUnreadyCount               *string   `json:"ok-total-unready-count,omitempty"`
	ScaleDownDelayAfterAdd            *string   `json:"scale-down-delay-after-add,omitempty"`
	ScaleDownDelayAfterDelete         *string   `json:"scale-down-delay-after-delete,omitempty"`
	ScaleDownDelayAfterFailure        *string   `json:"scale-down-delay-after-failure,omitempty"`
	ScaleDownUnneededTime             *string   `json:"scale-down-unneeded-time,omitempty"`
	ScaleDownUnreadyTime              *string   `json:"scale-down-unready-time,omitempty"`
	ScaleDownUtilizationThreshold     *string   `json:"scale-down-utilization-threshold,omitempty"`
	ScanInterval                      *string   `json:"scan-interval,omitempty"`
	SkipNodesWithLocalStorage         *string   `json:"skip-nodes-with-local-storage,omitempty"`
	SkipNodesWithSystemPods           *string   `json:"skip-nodes-with-system-pods,omitempty"`
}
