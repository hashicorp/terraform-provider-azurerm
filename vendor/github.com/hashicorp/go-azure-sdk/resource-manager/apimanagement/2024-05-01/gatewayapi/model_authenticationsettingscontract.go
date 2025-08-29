package gatewayapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthenticationSettingsContract struct {
	OAuth2                       *OAuth2AuthenticationSettingsContract   `json:"oAuth2,omitempty"`
	OAuth2AuthenticationSettings *[]OAuth2AuthenticationSettingsContract `json:"oAuth2AuthenticationSettings,omitempty"`
	Openid                       *OpenIdAuthenticationSettingsContract   `json:"openid,omitempty"`
	OpenidAuthenticationSettings *[]OpenIdAuthenticationSettingsContract `json:"openidAuthenticationSettings,omitempty"`
}
