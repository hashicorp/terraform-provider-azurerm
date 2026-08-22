package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureActiveDirectoryValidation struct {
	AllowedAudiences           *[]string                   `json:"allowedAudiences,omitempty"`
	DefaultAuthorizationPolicy *DefaultAuthorizationPolicy `json:"defaultAuthorizationPolicy,omitempty"`
	JwtClaimChecks             *JwtClaimChecks             `json:"jwtClaimChecks,omitempty"`
}
