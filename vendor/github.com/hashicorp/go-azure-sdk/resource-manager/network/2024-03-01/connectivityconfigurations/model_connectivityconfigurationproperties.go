package connectivityconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectivityConfigurationProperties struct {
	AppliesToGroups       []ConnectivityGroupItem `json:"appliesToGroups"`
	ConnectivityTopology  ConnectivityTopology    `json:"connectivityTopology"`
	DeleteExistingPeering *DeleteExistingPeering  `json:"deleteExistingPeering,omitempty"`
	Description           *string                 `json:"description,omitempty"`
	Hubs                  *[]Hub                  `json:"hubs,omitempty"`
	IsGlobal              *IsGlobal               `json:"isGlobal,omitempty"`
	ProvisioningState     *ProvisioningState      `json:"provisioningState,omitempty"`
	ResourceGuid          *string                 `json:"resourceGuid,omitempty"`
}
