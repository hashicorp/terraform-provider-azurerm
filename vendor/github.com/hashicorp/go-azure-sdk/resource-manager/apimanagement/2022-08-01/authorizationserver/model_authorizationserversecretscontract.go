package authorizationserver

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationServerSecretsContract struct {
	ClientSecret          *string `json:"clientSecret,omitempty"`
	ResourceOwnerPassword *string `json:"resourceOwnerPassword,omitempty"`
	ResourceOwnerUsername *string `json:"resourceOwnerUsername,omitempty"`
}
