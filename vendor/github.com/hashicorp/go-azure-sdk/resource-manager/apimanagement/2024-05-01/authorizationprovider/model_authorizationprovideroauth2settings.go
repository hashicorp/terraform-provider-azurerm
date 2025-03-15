package authorizationprovider

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationProviderOAuth2Settings struct {
	GrantTypes  *AuthorizationProviderOAuth2GrantTypes `json:"grantTypes,omitempty"`
	RedirectURL *string                                `json:"redirectUrl,omitempty"`
}
