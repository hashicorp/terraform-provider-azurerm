package partnernamespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartnerNamespaceProperties struct {
	DisableLocalAuth                    *bool                              `json:"disableLocalAuth,omitempty"`
	Endpoint                            *string                            `json:"endpoint,omitempty"`
	InboundIPRules                      *[]InboundIPRule                   `json:"inboundIpRules,omitempty"`
	PartnerRegistrationFullyQualifiedId *string                            `json:"partnerRegistrationFullyQualifiedId,omitempty"`
	PartnerTopicRoutingMode             *PartnerTopicRoutingMode           `json:"partnerTopicRoutingMode,omitempty"`
	PrivateEndpointConnections          *[]PrivateEndpointConnection       `json:"privateEndpointConnections,omitempty"`
	ProvisioningState                   *PartnerNamespaceProvisioningState `json:"provisioningState,omitempty"`
	PublicNetworkAccess                 *PublicNetworkAccess               `json:"publicNetworkAccess,omitempty"`
}
