package httprouteconfig

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HTTPRouteConfigProperties struct {
	CustomDomains      *[]CustomDomain                `json:"customDomains,omitempty"`
	Fqdn               *string                        `json:"fqdn,omitempty"`
	ProvisioningErrors *[]HTTPRouteProvisioningErrors `json:"provisioningErrors,omitempty"`
	ProvisioningState  *HTTPRouteProvisioningState    `json:"provisioningState,omitempty"`
	Rules              *[]HTTPRouteRule               `json:"rules,omitempty"`
}
