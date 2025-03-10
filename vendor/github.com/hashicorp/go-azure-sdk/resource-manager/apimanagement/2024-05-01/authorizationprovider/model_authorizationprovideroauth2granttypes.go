package authorizationprovider

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationProviderOAuth2GrantTypes struct {
	AuthorizationCode *map[string]string `json:"authorizationCode,omitempty"`
	ClientCredentials *map[string]string `json:"clientCredentials,omitempty"`
}
