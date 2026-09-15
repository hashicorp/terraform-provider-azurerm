package managedinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedInstancePrivateEndpointConnectionProperties struct {
	PrivateEndpoint                   *ManagedInstancePrivateEndpointProperty                   `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState *ManagedInstancePrivateLinkServiceConnectionStateProperty `json:"privateLinkServiceConnectionState,omitempty"`
	ProvisioningState                 *string                                                   `json:"provisioningState,omitempty"`
}
