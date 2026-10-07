package managedhsms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MHSMPrivateEndpointConnectionProperties struct {
	PrivateEndpoint                   *MHSMPrivateEndpoint                        `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState *MHSMPrivateLinkServiceConnectionState      `json:"privateLinkServiceConnectionState,omitempty"`
	ProvisioningState                 *PrivateEndpointConnectionProvisioningState `json:"provisioningState,omitempty"`
}
