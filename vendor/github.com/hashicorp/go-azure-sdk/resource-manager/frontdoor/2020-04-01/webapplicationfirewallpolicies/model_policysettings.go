package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicySettings struct {
	CustomBlockResponseBody       *string             `json:"customBlockResponseBody,omitempty"`
	CustomBlockResponseStatusCode *int64              `json:"customBlockResponseStatusCode,omitempty"`
	EnabledState                  *PolicyEnabledState `json:"enabledState,omitempty"`
	Mode                          *PolicyMode         `json:"mode,omitempty"`
	RedirectURL                   *string             `json:"redirectUrl,omitempty"`
}
