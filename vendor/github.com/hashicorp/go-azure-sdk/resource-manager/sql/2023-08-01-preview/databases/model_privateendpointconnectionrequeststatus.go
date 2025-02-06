package databases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointConnectionRequestStatus struct {
	PrivateEndpointConnectionName *string `json:"privateEndpointConnectionName,omitempty"`
	PrivateLinkServiceId          *string `json:"privateLinkServiceId,omitempty"`
	Status                        *string `json:"status,omitempty"`
}
