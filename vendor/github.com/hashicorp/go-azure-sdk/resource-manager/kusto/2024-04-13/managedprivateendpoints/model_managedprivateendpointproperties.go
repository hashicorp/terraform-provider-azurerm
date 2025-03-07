package managedprivateendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedPrivateEndpointProperties struct {
	GroupId                   string             `json:"groupId"`
	PrivateLinkResourceId     string             `json:"privateLinkResourceId"`
	PrivateLinkResourceRegion *string            `json:"privateLinkResourceRegion,omitempty"`
	ProvisioningState         *ProvisioningState `json:"provisioningState,omitempty"`
	RequestMessage            *string            `json:"requestMessage,omitempty"`
}
