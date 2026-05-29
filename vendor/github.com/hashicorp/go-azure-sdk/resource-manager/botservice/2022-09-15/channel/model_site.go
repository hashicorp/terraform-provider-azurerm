package channel

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Site struct {
	AppId                       *string   `json:"appId,omitempty"`
	ETag                        *string   `json:"eTag,omitempty"`
	IsBlockUserUploadEnabled    *bool     `json:"isBlockUserUploadEnabled,omitempty"`
	IsDetailedLoggingEnabled    *bool     `json:"isDetailedLoggingEnabled,omitempty"`
	IsEnabled                   bool      `json:"isEnabled"`
	IsEndpointParametersEnabled *bool     `json:"isEndpointParametersEnabled,omitempty"`
	IsNoStorageEnabled          *bool     `json:"isNoStorageEnabled,omitempty"`
	IsSecureSiteEnabled         *bool     `json:"isSecureSiteEnabled,omitempty"`
	IsTokenEnabled              *bool     `json:"isTokenEnabled,omitempty"`
	IsV1Enabled                 *bool     `json:"isV1Enabled,omitempty"`
	IsV3Enabled                 *bool     `json:"isV3Enabled,omitempty"`
	IsWebChatSpeechEnabled      *bool     `json:"isWebChatSpeechEnabled,omitempty"`
	IsWebchatPreviewEnabled     *bool     `json:"isWebchatPreviewEnabled,omitempty"`
	Key                         *string   `json:"key,omitempty"`
	Key2                        *string   `json:"key2,omitempty"`
	SiteId                      *string   `json:"siteId,omitempty"`
	SiteName                    string    `json:"siteName"`
	TenantId                    *string   `json:"tenantId,omitempty"`
	TrustedOrigins              *[]string `json:"trustedOrigins,omitempty"`
}
