package origins

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OriginUpdatePropertiesParameters struct {
	Enabled                    *bool   `json:"enabled,omitempty"`
	HTTPPort                   *int64  `json:"httpPort,omitempty"`
	HTTPSPort                  *int64  `json:"httpsPort,omitempty"`
	HostName                   *string `json:"hostName,omitempty"`
	OriginHostHeader           *string `json:"originHostHeader,omitempty"`
	Priority                   *int64  `json:"priority,omitempty"`
	PrivateLinkAlias           *string `json:"privateLinkAlias,omitempty"`
	PrivateLinkApprovalMessage *string `json:"privateLinkApprovalMessage,omitempty"`
	PrivateLinkLocation        *string `json:"privateLinkLocation,omitempty"`
	PrivateLinkResourceId      *string `json:"privateLinkResourceId,omitempty"`
	Weight                     *int64  `json:"weight,omitempty"`
}
