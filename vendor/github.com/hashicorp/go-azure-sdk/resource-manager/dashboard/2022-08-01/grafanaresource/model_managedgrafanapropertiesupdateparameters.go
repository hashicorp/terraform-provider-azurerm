package grafanaresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedGrafanaPropertiesUpdateParameters struct {
	ApiKey                  *ApiKey                  `json:"apiKey,omitempty"`
	DeterministicOutboundIP *DeterministicOutboundIP `json:"deterministicOutboundIP,omitempty"`
	GrafanaIntegrations     *GrafanaIntegrations     `json:"grafanaIntegrations,omitempty"`
	PublicNetworkAccess     *PublicNetworkAccess     `json:"publicNetworkAccess,omitempty"`
	ZoneRedundancy          *ZoneRedundancy          `json:"zoneRedundancy,omitempty"`
}
