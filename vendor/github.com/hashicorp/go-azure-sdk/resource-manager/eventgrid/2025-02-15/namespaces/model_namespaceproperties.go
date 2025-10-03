package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NamespaceProperties struct {
	InboundIPRules             *[]InboundIPRule             `json:"inboundIpRules,omitempty"`
	IsZoneRedundant            *bool                        `json:"isZoneRedundant,omitempty"`
	MinimumTlsVersionAllowed   *TlsVersion                  `json:"minimumTlsVersionAllowed,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	ProvisioningState          *NamespaceProvisioningState  `json:"provisioningState,omitempty"`
	PublicNetworkAccess        *PublicNetworkAccess         `json:"publicNetworkAccess,omitempty"`
	TopicSpacesConfiguration   *TopicSpacesConfiguration    `json:"topicSpacesConfiguration,omitempty"`
	TopicsConfiguration        *TopicsConfiguration         `json:"topicsConfiguration,omitempty"`
}
