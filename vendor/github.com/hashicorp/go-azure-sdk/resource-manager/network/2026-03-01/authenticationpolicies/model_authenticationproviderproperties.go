package authenticationpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthenticationProviderProperties struct {
	Audience          *string   `json:"audience,omitempty"`
	ClientId          string    `json:"clientId"`
	ClientSecret      *string   `json:"clientSecret,omitempty"`
	Issuer            string    `json:"issuer"`
	JwksUri           *string   `json:"jwksUri,omitempty"`
	Scope             *[]string `json:"scope,omitempty"`
	SessionCookieName *string   `json:"sessionCookieName,omitempty"`
	SessionTimeout    *string   `json:"sessionTimeout,omitempty"`
}
