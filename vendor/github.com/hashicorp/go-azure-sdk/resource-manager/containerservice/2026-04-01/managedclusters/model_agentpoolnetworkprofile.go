package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentPoolNetworkProfile struct {
	AllowedHostPorts          *[]PortRange `json:"allowedHostPorts,omitempty"`
	ApplicationSecurityGroups *[]string    `json:"applicationSecurityGroups,omitempty"`
	NodePublicIPTags          *[]IPTag     `json:"nodePublicIPTags,omitempty"`
}
