package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MachineUpdateProperties struct {
	AgentUpgrade               *AgentUpgrade  `json:"agentUpgrade,omitempty"`
	CloudMetadata              *CloudMetadata `json:"cloudMetadata,omitempty"`
	LocationData               *LocationData  `json:"locationData,omitempty"`
	OsProfile                  *OSProfile     `json:"osProfile,omitempty"`
	ParentClusterResourceId    *string        `json:"parentClusterResourceId,omitempty"`
	PrivateLinkScopeResourceId *string        `json:"privateLinkScopeResourceId,omitempty"`
}
