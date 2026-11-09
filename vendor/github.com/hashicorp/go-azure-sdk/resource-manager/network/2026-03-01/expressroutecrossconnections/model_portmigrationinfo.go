package expressroutecrossconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PortMigrationInfo struct {
	FailureReason   *string          `json:"failureReason,omitempty"`
	Peerings        *[]PeeringHealth `json:"peerings,omitempty"`
	Phase           *string          `json:"phase,omitempty"`
	PortId          *string          `json:"portId,omitempty"`
	SourcePortId    *string          `json:"sourcePortId,omitempty"`
	SourcePortStats *SourcePortStats `json:"sourcePortStats,omitempty"`
	Status          *string          `json:"status,omitempty"`
}
