package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DevToolPortalProperties struct {
	Components        *[]DevToolPortalComponent       `json:"components,omitempty"`
	Features          *DevToolPortalFeatureSettings   `json:"features,omitempty"`
	ProvisioningState *DevToolPortalProvisioningState `json:"provisioningState,omitempty"`
	Public            *bool                           `json:"public,omitempty"`
	SsoProperties     *DevToolPortalSsoProperties     `json:"ssoProperties,omitempty"`
	Url               *string                         `json:"url,omitempty"`
}
