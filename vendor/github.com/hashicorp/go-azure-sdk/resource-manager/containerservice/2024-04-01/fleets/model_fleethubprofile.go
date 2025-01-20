package fleets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FleetHubProfile struct {
	AgentProfile           *AgentProfile           `json:"agentProfile,omitempty"`
	ApiServerAccessProfile *APIServerAccessProfile `json:"apiServerAccessProfile,omitempty"`
	DnsPrefix              *string                 `json:"dnsPrefix,omitempty"`
	Fqdn                   *string                 `json:"fqdn,omitempty"`
	KubernetesVersion      *string                 `json:"kubernetesVersion,omitempty"`
	PortalFqdn             *string                 `json:"portalFqdn,omitempty"`
}
