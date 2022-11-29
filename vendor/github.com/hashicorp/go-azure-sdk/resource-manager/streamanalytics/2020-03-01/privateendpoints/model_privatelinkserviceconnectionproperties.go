package privateendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkServiceConnectionProperties struct {
	GroupIds                          *[]string                   `json:"groupIds,omitempty"`
	PrivateLinkServiceConnectionState *PrivateLinkConnectionState `json:"privateLinkServiceConnectionState,omitempty"`
	PrivateLinkServiceId              *string                     `json:"privateLinkServiceId,omitempty"`
	RequestMessage                    *string                     `json:"requestMessage,omitempty"`
}
