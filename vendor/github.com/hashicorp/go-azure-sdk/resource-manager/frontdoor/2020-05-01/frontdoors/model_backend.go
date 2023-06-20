package frontdoors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Backend struct {
	Address                    *string                `json:"address,omitempty"`
	BackendHostHeader          *string                `json:"backendHostHeader,omitempty"`
	EnabledState               *BackendEnabledState   `json:"enabledState,omitempty"`
	HTTPPort                   *int64                 `json:"httpPort,omitempty"`
	HTTPSPort                  *int64                 `json:"httpsPort,omitempty"`
	Priority                   *int64                 `json:"priority,omitempty"`
	PrivateEndpointStatus      *PrivateEndpointStatus `json:"privateEndpointStatus,omitempty"`
	PrivateLinkAlias           *string                `json:"privateLinkAlias,omitempty"`
	PrivateLinkApprovalMessage *string                `json:"privateLinkApprovalMessage,omitempty"`
	PrivateLinkLocation        *string                `json:"privateLinkLocation,omitempty"`
	PrivateLinkResourceId      *string                `json:"privateLinkResourceId,omitempty"`
	Weight                     *int64                 `json:"weight,omitempty"`
}
