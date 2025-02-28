package authorizationprovider

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationProviderContractProperties struct {
	DisplayName      *string                              `json:"displayName,omitempty"`
	IdentityProvider *string                              `json:"identityProvider,omitempty"`
	Oauth2           *AuthorizationProviderOAuth2Settings `json:"oauth2,omitempty"`
}
