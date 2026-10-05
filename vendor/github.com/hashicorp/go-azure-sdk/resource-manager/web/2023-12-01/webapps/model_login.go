package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Login struct {
	AllowedExternalRedirectURLs   *[]string         `json:"allowedExternalRedirectUrls,omitempty"`
	CookieExpiration              *CookieExpiration `json:"cookieExpiration,omitempty"`
	Nonce                         *Nonce            `json:"nonce,omitempty"`
	PreserveURLFragmentsForLogins *bool             `json:"preserveUrlFragmentsForLogins,omitempty"`
	Routes                        *LoginRoutes      `json:"routes,omitempty"`
	TokenStore                    *TokenStore       `json:"tokenStore,omitempty"`
}
