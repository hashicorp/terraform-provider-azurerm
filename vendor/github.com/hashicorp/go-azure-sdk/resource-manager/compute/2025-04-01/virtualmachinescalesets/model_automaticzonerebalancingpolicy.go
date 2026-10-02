package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomaticZoneRebalancingPolicy struct {
	Enabled           *bool              `json:"enabled,omitempty"`
	RebalanceBehavior *RebalanceBehavior `json:"rebalanceBehavior,omitempty"`
	RebalanceStrategy *RebalanceStrategy `json:"rebalanceStrategy,omitempty"`
}
