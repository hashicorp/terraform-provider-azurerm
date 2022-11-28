package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerPrivateEndpointConnectionProperties struct {
	PrivateEndpoint                   *PrivateEndpointProperty                         `json:"privateEndpoint"`
	PrivateLinkServiceConnectionState *ServerPrivateLinkServiceConnectionStateProperty `json:"privateLinkServiceConnectionState"`
	ProvisioningState                 *PrivateEndpointProvisioningState                `json:"provisioningState,omitempty"`
}
