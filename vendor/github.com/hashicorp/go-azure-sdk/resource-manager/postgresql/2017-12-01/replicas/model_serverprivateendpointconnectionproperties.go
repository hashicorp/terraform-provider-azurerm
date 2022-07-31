package replicas

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerPrivateEndpointConnectionProperties struct {
	PrivateEndpoint                   *PrivateEndpointProperty                         `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState *ServerPrivateLinkServiceConnectionStateProperty `json:"privateLinkServiceConnectionState,omitempty"`
	ProvisioningState                 *PrivateEndpointProvisioningState                `json:"provisioningState,omitempty"`
}
