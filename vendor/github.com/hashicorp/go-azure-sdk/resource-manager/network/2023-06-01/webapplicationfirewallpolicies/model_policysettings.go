package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicySettings struct {
	CustomBlockResponseBody       *string                             `json:"customBlockResponseBody,omitempty"`
	CustomBlockResponseStatusCode *int64                              `json:"customBlockResponseStatusCode,omitempty"`
	FileUploadEnforcement         *bool                               `json:"fileUploadEnforcement,omitempty"`
	FileUploadLimitInMb           *int64                              `json:"fileUploadLimitInMb,omitempty"`
	LogScrubbing                  *PolicySettingsLogScrubbing         `json:"logScrubbing,omitempty"`
	MaxRequestBodySizeInKb        *int64                              `json:"maxRequestBodySizeInKb,omitempty"`
	Mode                          *WebApplicationFirewallMode         `json:"mode,omitempty"`
	RequestBodyCheck              *bool                               `json:"requestBodyCheck,omitempty"`
	RequestBodyEnforcement        *bool                               `json:"requestBodyEnforcement,omitempty"`
	RequestBodyInspectLimitInKB   *int64                              `json:"requestBodyInspectLimitInKB,omitempty"`
	State                         *WebApplicationFirewallEnabledState `json:"state,omitempty"`
}
