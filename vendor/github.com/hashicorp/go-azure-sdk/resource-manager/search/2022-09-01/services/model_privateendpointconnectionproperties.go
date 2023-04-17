package services

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointConnectionProperties struct {
	GroupId                           *string                                                               `json:"groupId,omitempty"`
	PrivateEndpoint                   *PrivateEndpointConnectionPropertiesPrivateEndpoint                   `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState *PrivateEndpointConnectionPropertiesPrivateLinkServiceConnectionState `json:"privateLinkServiceConnectionState,omitempty"`
	ProvisioningState                 *PrivateLinkServiceConnectionProvisioningState                        `json:"provisioningState,omitempty"`
}
