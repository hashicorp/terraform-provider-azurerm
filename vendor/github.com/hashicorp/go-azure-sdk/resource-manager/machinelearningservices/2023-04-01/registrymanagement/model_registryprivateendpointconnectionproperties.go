package registrymanagement

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryPrivateEndpointConnectionProperties struct {
	GroupIds                          *[]string                                  `json:"groupIds,omitempty"`
	PrivateEndpoint                   *PrivateEndpointResource                   `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState *RegistryPrivateLinkServiceConnectionState `json:"privateLinkServiceConnectionState,omitempty"`
	ProvisioningState                 *string                                    `json:"provisioningState,omitempty"`
}
