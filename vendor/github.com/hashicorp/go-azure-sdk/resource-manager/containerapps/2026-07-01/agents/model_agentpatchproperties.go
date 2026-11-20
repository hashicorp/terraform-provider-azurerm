package agents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentPatchProperties struct {
	ActionConfiguration             *ActionConfiguration             `json:"actionConfiguration,omitempty"`
	AgentIdentity                   *AgentIdentityPatch              `json:"agentIdentity,omitempty"`
	DefaultModel                    *DefaultModel                    `json:"defaultModel,omitempty"`
	IncidentManagementConfiguration *IncidentManagementConfiguration `json:"incidentManagementConfiguration,omitempty"`
	KnowledgeGraphConfiguration     *KnowledgeGraphConfiguration     `json:"knowledgeGraphConfiguration,omitempty"`
	LogConfiguration                *LogConfiguration                `json:"logConfiguration,omitempty"`
	SandboxConfiguration            *SandboxConfiguration            `json:"sandboxConfiguration,omitempty"`
	UpgradeChannel                  *UpgradeChannel                  `json:"upgradeChannel,omitempty"`
	VnetConfiguration               *VnetConfiguration               `json:"vnetConfiguration,omitempty"`
}
