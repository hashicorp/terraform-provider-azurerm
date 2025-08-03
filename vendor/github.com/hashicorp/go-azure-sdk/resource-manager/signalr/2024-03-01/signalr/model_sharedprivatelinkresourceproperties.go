package signalr

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharedPrivateLinkResourceProperties struct {
	GroupId               string                           `json:"groupId"`
	PrivateLinkResourceId string                           `json:"privateLinkResourceId"`
	ProvisioningState     *ProvisioningState               `json:"provisioningState,omitempty"`
	RequestMessage        *string                          `json:"requestMessage,omitempty"`
	Status                *SharedPrivateLinkResourceStatus `json:"status,omitempty"`
}
