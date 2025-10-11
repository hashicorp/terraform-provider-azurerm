package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopologyParameters struct {
	TargetResourceGroupName *string      `json:"targetResourceGroupName,omitempty"`
	TargetSubnet            *SubResource `json:"targetSubnet,omitempty"`
	TargetVirtualNetwork    *SubResource `json:"targetVirtualNetwork,omitempty"`
}
