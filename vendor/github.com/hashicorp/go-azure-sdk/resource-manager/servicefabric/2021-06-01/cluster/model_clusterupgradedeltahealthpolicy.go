package cluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterUpgradeDeltaHealthPolicy struct {
	ApplicationDeltaHealthPolicies             *map[string]ApplicationDeltaHealthPolicy `json:"applicationDeltaHealthPolicies,omitempty"`
	MaxPercentDeltaUnhealthyApplications       int64                                    `json:"maxPercentDeltaUnhealthyApplications"`
	MaxPercentDeltaUnhealthyNodes              int64                                    `json:"maxPercentDeltaUnhealthyNodes"`
	MaxPercentUpgradeDomainDeltaUnhealthyNodes int64                                    `json:"maxPercentUpgradeDomainDeltaUnhealthyNodes"`
}
