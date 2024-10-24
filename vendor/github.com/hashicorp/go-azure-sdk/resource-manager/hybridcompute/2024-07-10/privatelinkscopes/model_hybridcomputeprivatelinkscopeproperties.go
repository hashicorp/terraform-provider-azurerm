package privatelinkscopes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HybridComputePrivateLinkScopeProperties struct {
	PrivateEndpointConnections *[]PrivateEndpointConnectionDataModel `json:"privateEndpointConnections,omitempty"`
	PrivateLinkScopeId         *string                               `json:"privateLinkScopeId,omitempty"`
	ProvisioningState          *string                               `json:"provisioningState,omitempty"`
	PublicNetworkAccess        *PublicNetworkAccessType              `json:"publicNetworkAccess,omitempty"`
}
