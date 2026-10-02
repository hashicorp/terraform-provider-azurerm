package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OpenIdConnectRegistration struct {
	ClientCredential           *OpenIdConnectClientCredential `json:"clientCredential,omitempty"`
	ClientId                   *string                        `json:"clientId,omitempty"`
	OpenIdConnectConfiguration *OpenIdConnectConfig           `json:"openIdConnectConfiguration,omitempty"`
}
