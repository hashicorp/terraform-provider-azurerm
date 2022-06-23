package apps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppProperties struct {
	ApplicationId              *string                      `json:"applicationId,omitempty"`
	DisplayName                *string                      `json:"displayName,omitempty"`
	NetworkRuleSets            *NetworkRuleSets             `json:"networkRuleSets,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	ProvisioningState          *ProvisioningState           `json:"provisioningState,omitempty"`
	PublicNetworkAccess        *PublicNetworkAccess         `json:"publicNetworkAccess,omitempty"`
	State                      *AppState                    `json:"state,omitempty"`
	Subdomain                  *string                      `json:"subdomain,omitempty"`
	Template                   *string                      `json:"template,omitempty"`
}
