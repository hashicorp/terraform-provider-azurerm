package dicomservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointConnectionProperties struct {
	PrivateEndpoint                   *PrivateEndpoint                            `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState PrivateLinkServiceConnectionState           `json:"privateLinkServiceConnectionState"`
	ProvisioningState                 *PrivateEndpointConnectionProvisioningState `json:"provisioningState,omitempty"`
}
