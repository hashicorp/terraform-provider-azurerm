package apimanagementservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointConnectionWrapperProperties struct {
	GroupIds                          *[]string                         `json:"groupIds,omitempty"`
	PrivateEndpoint                   *ArmIdWrapper                     `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState PrivateLinkServiceConnectionState `json:"privateLinkServiceConnectionState"`
	ProvisioningState                 *string                           `json:"provisioningState,omitempty"`
}
