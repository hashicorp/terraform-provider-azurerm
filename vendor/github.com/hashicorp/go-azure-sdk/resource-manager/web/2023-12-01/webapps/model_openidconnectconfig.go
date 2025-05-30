package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OpenIdConnectConfig struct {
	AuthorizationEndpoint        *string `json:"authorizationEndpoint,omitempty"`
	CertificationUri             *string `json:"certificationUri,omitempty"`
	Issuer                       *string `json:"issuer,omitempty"`
	TokenEndpoint                *string `json:"tokenEndpoint,omitempty"`
	WellKnownOpenIdConfiguration *string `json:"wellKnownOpenIdConfiguration,omitempty"`
}
