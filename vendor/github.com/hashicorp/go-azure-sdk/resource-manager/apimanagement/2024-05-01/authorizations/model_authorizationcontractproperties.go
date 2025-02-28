package authorizations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationContractProperties struct {
	AuthorizationType *AuthorizationType  `json:"authorizationType,omitempty"`
	Error             *AuthorizationError `json:"error,omitempty"`
	Oauth2grantType   *OAuth2GrantType    `json:"oauth2grantType,omitempty"`
	Parameters        *map[string]string  `json:"parameters,omitempty"`
	Status            *string             `json:"status,omitempty"`
}
