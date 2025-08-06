package managedprivateendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedPrivateEndpointModelProperties struct {
	ConnectionState             *ManagedPrivateEndpointConnectionState `json:"connectionState,omitempty"`
	GroupIds                    *[]string                              `json:"groupIds,omitempty"`
	PrivateLinkResourceId       *string                                `json:"privateLinkResourceId,omitempty"`
	PrivateLinkResourceRegion   *string                                `json:"privateLinkResourceRegion,omitempty"`
	PrivateLinkServicePrivateIP *string                                `json:"privateLinkServicePrivateIP,omitempty"`
	PrivateLinkServiceURL       *string                                `json:"privateLinkServiceUrl,omitempty"`
	ProvisioningState           *ProvisioningState                     `json:"provisioningState,omitempty"`
	RequestMessage              *string                                `json:"requestMessage,omitempty"`
}
