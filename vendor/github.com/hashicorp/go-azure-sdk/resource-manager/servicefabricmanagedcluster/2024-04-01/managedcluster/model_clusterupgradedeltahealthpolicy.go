package managedcluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterUpgradeDeltaHealthPolicy struct {
	MaxPercentDeltaUnhealthyApplications       *int64 `json:"maxPercentDeltaUnhealthyApplications,omitempty"`
	MaxPercentDeltaUnhealthyNodes              int64  `json:"maxPercentDeltaUnhealthyNodes"`
	MaxPercentUpgradeDomainDeltaUnhealthyNodes *int64 `json:"maxPercentUpgradeDomainDeltaUnhealthyNodes,omitempty"`
}
