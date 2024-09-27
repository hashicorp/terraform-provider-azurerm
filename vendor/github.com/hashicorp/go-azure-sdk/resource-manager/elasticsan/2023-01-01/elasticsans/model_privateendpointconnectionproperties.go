package elasticsans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointConnectionProperties struct {
	GroupIds                          *[]string                         `json:"groupIds,omitempty"`
	PrivateEndpoint                   *PrivateEndpoint                  `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState PrivateLinkServiceConnectionState `json:"privateLinkServiceConnectionState"`
	ProvisioningState                 *ProvisioningStates               `json:"provisioningState,omitempty"`
}
