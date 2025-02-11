package connectedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArcAgentProfile struct {
	AgentAutoUpgrade    *AutoUpgradeOptions `json:"agentAutoUpgrade,omitempty"`
	DesiredAgentVersion *string             `json:"desiredAgentVersion,omitempty"`
}
