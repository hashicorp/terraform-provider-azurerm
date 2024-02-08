package agentpools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentPoolUpgradeProfileProperties struct {
	KubernetesVersion      string                                              `json:"kubernetesVersion"`
	LatestNodeImageVersion *string                                             `json:"latestNodeImageVersion,omitempty"`
	OsType                 OSType                                              `json:"osType"`
	Upgrades               *[]AgentPoolUpgradeProfilePropertiesUpgradesInlined `json:"upgrades,omitempty"`
}
