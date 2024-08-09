package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointConnectionProperties struct {
	GroupIds                          *[]string                                   `json:"groupIds,omitempty"`
	PrivateEndpoint                   *RemotePrivateEndpointConnection            `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState *RemotePrivateLinkServiceConnectionState    `json:"privateLinkServiceConnectionState,omitempty"`
	ProvisioningState                 *PrivateEndpointConnectionProvisioningState `json:"provisioningState,omitempty"`
}
