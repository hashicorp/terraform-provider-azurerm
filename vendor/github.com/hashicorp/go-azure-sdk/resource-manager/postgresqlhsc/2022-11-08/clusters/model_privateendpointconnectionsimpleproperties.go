package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointConnectionSimpleProperties struct {
	GroupIds                          *[]string                          `json:"groupIds,omitempty"`
	PrivateEndpoint                   *PrivateEndpointProperty           `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState *PrivateLinkServiceConnectionState `json:"privateLinkServiceConnectionState,omitempty"`
}
