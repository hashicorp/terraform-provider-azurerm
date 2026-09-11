package publicipaddresses

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkServiceConnectionProperties struct {
	GroupIds                          *[]string                          `json:"groupIds,omitempty"`
	PrivateLinkServiceConnectionState *PrivateLinkServiceConnectionState `json:"privateLinkServiceConnectionState,omitempty"`
	PrivateLinkServiceId              *string                            `json:"privateLinkServiceId,omitempty"`
	ProvisioningState                 *ProvisioningState                 `json:"provisioningState,omitempty"`
	RequestMessage                    *string                            `json:"requestMessage,omitempty"`
}
