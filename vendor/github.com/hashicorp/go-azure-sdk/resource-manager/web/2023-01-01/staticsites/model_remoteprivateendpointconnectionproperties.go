package staticsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemotePrivateEndpointConnectionProperties struct {
	IPAddresses                       *[]string                   `json:"ipAddresses,omitempty"`
	PrivateEndpoint                   *ArmIdWrapper               `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState *PrivateLinkConnectionState `json:"privateLinkServiceConnectionState,omitempty"`
	ProvisioningState                 *string                     `json:"provisioningState,omitempty"`
}
