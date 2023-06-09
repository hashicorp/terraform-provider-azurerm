package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicySettings struct {
	CustomBlockResponseBody       *string                             `json:"customBlockResponseBody,omitempty"`
	CustomBlockResponseStatusCode *int64                              `json:"customBlockResponseStatusCode,omitempty"`
	FileUploadLimitInMb           *int64                              `json:"fileUploadLimitInMb,omitempty"`
	MaxRequestBodySizeInKb        *int64                              `json:"maxRequestBodySizeInKb,omitempty"`
	Mode                          *WebApplicationFirewallMode         `json:"mode,omitempty"`
	RequestBodyCheck              *bool                               `json:"requestBodyCheck,omitempty"`
	State                         *WebApplicationFirewallEnabledState `json:"state,omitempty"`
}
