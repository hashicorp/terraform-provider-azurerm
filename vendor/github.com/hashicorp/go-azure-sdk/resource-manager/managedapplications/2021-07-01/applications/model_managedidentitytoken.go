package applications

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedIdentityToken struct {
	AccessToken           *string `json:"accessToken,omitempty"`
	AuthorizationAudience *string `json:"authorizationAudience,omitempty"`
	ExpiresIn             *string `json:"expiresIn,omitempty"`
	ExpiresOn             *string `json:"expiresOn,omitempty"`
	NotBefore             *string `json:"notBefore,omitempty"`
	ResourceId            *string `json:"resourceId,omitempty"`
	TokenType             *string `json:"tokenType,omitempty"`
}
