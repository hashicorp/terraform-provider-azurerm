package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointConnectionProperties struct {
	GroupId                           *string                                    `json:"groupId,omitempty"`
	PrivateEndpoint                   *PrivateEndpointProperty                   `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState *PrivateLinkServiceConnectionStateProperty `json:"privateLinkServiceConnectionState,omitempty"`
	ProvisioningState                 *string                                    `json:"provisioningState,omitempty"`
}
