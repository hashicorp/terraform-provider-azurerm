package managedapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiOAuthSettings struct {
	ClientId         *string                               `json:"clientId,omitempty"`
	ClientSecret     *string                               `json:"clientSecret,omitempty"`
	CustomParameters *map[string]ApiOAuthSettingsParameter `json:"customParameters,omitempty"`
	IdentityProvider *string                               `json:"identityProvider,omitempty"`
	Properties       *interface{}                          `json:"properties,omitempty"`
	RedirectUrl      *string                               `json:"redirectUrl,omitempty"`
	Scopes           *[]string                             `json:"scopes,omitempty"`
}
