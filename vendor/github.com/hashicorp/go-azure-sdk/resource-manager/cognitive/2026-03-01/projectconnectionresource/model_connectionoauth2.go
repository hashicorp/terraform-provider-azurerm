package projectconnectionresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionOAuth2 struct {
	AuthURL        *string `json:"authUrl,omitempty"`
	ClientId       *string `json:"clientId,omitempty"`
	ClientSecret   *string `json:"clientSecret,omitempty"`
	DeveloperToken *string `json:"developerToken,omitempty"`
	Password       *string `json:"password,omitempty"`
	RefreshToken   *string `json:"refreshToken,omitempty"`
	TenantId       *string `json:"tenantId,omitempty"`
	Username       *string `json:"username,omitempty"`
}
